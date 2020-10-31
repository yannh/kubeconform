package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"github.com/yannh/kubeconform/pkg/config"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/fsutils"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

type validationResult struct {
	filename, kind, version, Name string
	err                           error
	skipped                       bool
}

func resourcesFromReader(r io.Reader) ([][]byte, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return [][]byte{}, err
	}

	resources := bytes.Split(data, []byte("---\n"))

	return resources, nil
}

func downloadSchema(registries []registry.Registry, kind, version, k8sVersion string) (*gojsonschema.Schema, error) {
	var err error
	var schemaBytes []byte

	for _, reg := range registries {
		schemaBytes, err = reg.DownloadSchema(kind, version, k8sVersion)
		if err == nil {
			return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaBytes))
		}

		// If we get a 404, we try the next registry, but we exit if we get a real failure
		if er, retryable := err.(registry.Retryable); retryable && !er.IsRetryable() {
			continue
		}

		return nil, err
	}

	return nil, nil // No schema found - we don't consider it an error, resource will be skipped
}

// filter returns true if the file should be skipped
// Returning an array, this Reader might container multiple resources
func ValidateStream(r io.Reader, regs []registry.Registry, k8sVersion string, c *cache.SchemaCache, skip func(signature resource.Signature) bool, ignoreMissingSchemas bool) []validationResult {
	rawResources, err := resourcesFromReader(r)
	if err != nil {
		return []validationResult{{err: fmt.Errorf("failed reading file: %s", err)}}
	}

	validationResults := []validationResult{}
	if len(rawResources) == 0 {
		// In case a file has no resources, we want to capture that the file was parsed - and therefore send a message with an empty resource and no error
		validationResults = append(validationResults, validationResult{kind: "", version: "", Name: "", err: nil, skipped: false})
	}

	for _, rawResource := range rawResources {
		var sig resource.Signature
		if sig, err = resource.SignatureFromBytes(rawResource); err != nil {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, Name: sig.Name, err: fmt.Errorf("error while parsing: %s", err)})
			continue
		}

		if sig.Kind == "" {
			validationResults = append(validationResults, validationResult{kind: "", version: "", Name: "", err: nil, skipped: false})
			continue // We skip resoures that don't have a Kind defined
		}

		if skip(sig) {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, Name: sig.Name, err: nil, skipped: true})
			continue
		}

		ok := false
		var schema *gojsonschema.Schema
		cacheKey := ""

		if c != nil {
			cacheKey = cache.Key(sig.Kind, sig.Version, k8sVersion)
			schema, ok = c.Get(cacheKey)
		}

		if !ok {
			schema, err = downloadSchema(regs, sig.Kind, sig.Version, k8sVersion)
			if err != nil {
				validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, Name: sig.Name, err: err, skipped: false})
				continue
			}

			if c != nil {
				c.Set(cacheKey, schema)
			}
		}

		if schema == nil {
			if ignoreMissingSchemas {
				validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, Name: sig.Name, err: nil, skipped: true})
			} else {
				validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, Name: sig.Name, err: fmt.Errorf("could not find schema for %s", sig.Kind), skipped: false})
			}
		}

		err = validator.Validate(rawResource, schema)
		validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, Name: sig.Name, err: err})
	}

	return validationResults
}

func processResults(o output.Output, validationResults chan []validationResult, result chan<- bool) {
	success := true
	for results := range validationResults {
		for _, result := range results {
			if result.err != nil {
				success = false
			}

			if err := o.Write(result.filename, result.kind, result.Name, result.version, result.err, result.skipped); err != nil {
				fmt.Fprint(os.Stderr, "failed writing log\n")
			}
		}
	}

	result <- success
}

func getFiles(files []string, fileBatches chan []string, validationResults chan []validationResult) {
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			validationResults <- []validationResult{{
				filename: filename,
				err:      err,
				skipped:  false,
			}}
			continue
		}
		defer file.Close()

		fi, err := file.Stat()
		switch {
		case err != nil:
			validationResults <- []validationResult{{
				filename: filename,
				err:      err,
				skipped:  false,
			}}

		case fi.IsDir():
			if err := fsutils.FindYamlInDir(filename, fileBatches, 10); err != nil {
				validationResults <- []validationResult{{
					filename: filename,
					err:      err,
					skipped:  false,
				}}
			}

		default:
			fileBatches <- []string{filename}
		}
	}
}

func realMain() int {
	var err error

	cfg := config.FromFlags()

	if cfg.Help {
		flag.Usage()
		return 1
	}

	// Detect whether we have data being piped through stdin
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if len(flag.Args()) == 1 && flag.Args()[0] == "-" {
		isStdin = true
	}

	filter := func(signature resource.Signature) bool {
		isSkipKind, ok := cfg.SkipKinds[signature.Kind]
		return ok && isSkipKind
	}

	registries := []registry.Registry{}
	for _, schemaLocation := range cfg.SchemaLocations {
		if !strings.HasSuffix(schemaLocation, "json") { // If we dont specify a full templated path, we assume the paths of kubernetesjsonschema.dev
			schemaLocation += "/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json"
		}

		if strings.HasPrefix(schemaLocation, "http") {
			registries = append(registries, registry.NewHTTPRegistry(schemaLocation, cfg.Strict))
		} else {
			registries = append(registries, registry.NewLocalRegistry(schemaLocation, cfg.Strict))
		}
	}

	validationResults := make(chan []validationResult)
	c := cache.New()

	fileBatches := make(chan []string)
	go func() {
		getFiles(cfg.Files, fileBatches, validationResults)
		close(fileBatches)
	}()

	var o output.Output
	if o, err = output.New(cfg.OutputFormat, cfg.Summary, isStdin, cfg.Verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	res := make(chan bool)
	go processResults(o, validationResults, res)

	if isStdin {
		res := ValidateStream(os.Stdin, registries, cfg.KubernetesVersion, c, filter, cfg.IgnoreMissingSchemas)
		for i := range res {
			res[i].filename = "stdin"
		}
		validationResults <- res
	} else {
		var wg sync.WaitGroup
		for i := 0; i < cfg.NumberOfWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for fileBatch := range fileBatches {
					for _, filename := range fileBatch {
						f, err := os.Open(filename)
						if err != nil {
							validationResults <- []validationResult{{
								filename: filename,
								err:      err,
								skipped:  true,
							}}
							continue
						}

						res := ValidateStream(f, registries, cfg.KubernetesVersion, c, filter, cfg.IgnoreMissingSchemas)
						f.Close()

						for i := range res {
							res[i].filename = filename
						}
						validationResults <- res
					}
				}
			}()
		}

		wg.Wait()
	}

	close(validationResults)
	success := <-res
	o.Flush()

	if !success {
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
