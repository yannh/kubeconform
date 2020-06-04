package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"github.com/yannh/kubeconform/pkg/output"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/fsutils"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

type validationResult struct {
	filename, kind, version string
	err                     error
	skipped                 bool
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
	var schema *gojsonschema.Schema
	var schemaBytes []byte

	for _, reg := range registries {
		schemaBytes, err = reg.DownloadSchema(kind, version, k8sVersion)

		if err != nil {
			// If we get a 404, we keep trying, but we exit if we get a real failure
			if er, retryable := err.(registry.Retryable); !retryable || er.IsRetryable() {
				return nil, err
			}

			continue // 404 from this registry, try next registry
		}

		if schema, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaBytes)); err != nil {
			return nil, err // Got a schema, but fail to parse it
		}

		return schema, nil
	}

	return nil, nil // No schema found - we don't consider it an error, resource willb e skipped
}

// filter returns true if the file should be skipped
// Returning an array, this Reader might container multiple resources
func ValidateStream(r io.Reader, regs []registry.Registry, k8sVersion string, c *cache.SchemaCache, skip func(signature resource.Signature) bool, ignoreMissingSchemas bool) []validationResult {
	rawResources, err := resourcesFromReader(r)
	if err != nil {
		return []validationResult{{err: fmt.Errorf("failed reading file: %s", err)}}
	}

	validationResults := []validationResult{}

	for _, rawResource := range rawResources {
		var sig resource.Signature
		if sig, err = resource.SignatureFromBytes(rawResource); err != nil {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: fmt.Errorf("error while parsing: %s", err)})
			continue
		}

		if sig.Kind == "" {
			continue // We skip resoures that don't have a Kind defined
		}

		if skip(sig) {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil, skipped: true})
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
				validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: err, skipped: false})
				continue
			} else if schema == nil { // skip if no schema was found, but there was no error TODO: Fail by default, provide a -skip-missing-schema
				if ignoreMissingSchemas {
					validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil, skipped: true})
				} else {
					validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: fmt.Errorf("could not find schema for %s", sig.Kind), skipped: false})
				}
				if c != nil {
					c.Set(cacheKey, nil)
				}
				continue
			} else if c != nil {
				c.Set(cacheKey, schema)
			}
		}

		err = validator.Validate(rawResource, schema)
		validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: err})
	}

	return validationResults
}

type arrayParam []string

func (ap *arrayParam) String() string {
	return strings.Join(*ap, " - ")
}

func (ap *arrayParam) Set(value string) error {
	*ap = append(*ap, value)
	return nil
}

func getLogger(outputFormat string, printSummary, verbose bool) (output.Output, error) {
	w := os.Stdout

	switch {
	case outputFormat == "text":
		return output.Text(w, printSummary, verbose), nil
	case outputFormat == "json":
		return output.JSON(w, printSummary, verbose), nil
	default:
		return nil, fmt.Errorf("-output must be text or json")
	}
}

func skipKindsMap(skipKindsCSV string) map[string]bool {
	splitKinds := strings.Split(skipKindsCSV, ",")
	skipKinds := map[string]bool{}
	for _, kind := range splitKinds {
		if len(kind) > 0 {
			skipKinds[kind] = true
		}
	}
	return skipKinds
}

func processResults(o output.Output, validationResults chan []validationResult, result chan<- bool) {
	success := true
	for results := range validationResults {
		for _, result := range results {
			if result.err != nil {
				success = false
			}

			if err := o.Write(result.filename, result.kind, result.version, result.err, result.skipped); err != nil {
				fmt.Fprint(os.Stderr, "failed writing log\n")
			}
		}
	}

	result <- success
}

func realMain() int {
	var files, dirs, schemas arrayParam
	var skipKindsCSV, k8sVersion, outputFormat string
	var summary, strict, verbose, ignoreMissingSchemas bool
	var nWorkers int
	var err error

	flag.StringVar(&k8sVersion, "k8sversion", "1.18.0", "version of Kubernetes to test against")
	flag.Var(&files, "file", "file to validate (can be specified multiple times)")
	flag.Var(&dirs, "dir", "directory to validate (can be specified multiple times)")
	flag.Var(&schemas, "schema", "file containing an additional Schema (can be specified multiple times)")
	flag.BoolVar(&ignoreMissingSchemas, "ignore-missing-schemas", false, "skip files with missing schemas instead of failing")
	flag.BoolVar(&summary, "summary", false, "print a summary at the end")
	flag.IntVar(&nWorkers, "n", 4, "number of routines to run in parallel")
	flag.StringVar(&skipKindsCSV, "skip", "", "comma-separated list of kinds to ignore")
	flag.BoolVar(&strict, "strict", false, "disallow additional properties not in schema")
	flag.StringVar(&outputFormat, "output", "text", "output format - text, json")
	flag.BoolVar(&verbose, "verbose", false, "print results for all resources")
	flag.Parse()

	skipKinds := skipKindsMap(skipKindsCSV)

	filter := func(signature resource.Signature) bool {
		isSkipKind, ok := skipKinds[signature.Kind]
		return ok && isSkipKind
	}

	registries := []registry.Registry{}
	registries = append(registries, registry.NewKubernetesRegistry(strict))
	if len(schemas) > 0 {
		localRegistry, err := registry.NewLocalSchemas(schemas)
		if err != nil {
			log.Fatalf("%s", err)
		}
		registries = append(registries, localRegistry)
	}

	fileBatches := make(chan []string)
	go func() {
		for _, dir := range dirs {
			if err := fsutils.FindYamlInDir(dir, fileBatches, 10); err != nil {
				log.Printf("failed processing folder %s: %s", dir, err)
			}
		}

		for _, filename := range files {
			fileBatches <- []string{filename}
		}

		close(fileBatches)
	}()

	var o output.Output
	if o, err = getLogger(outputFormat, summary, verbose); err != nil {
		fmt.Println(err)
		return 1
	}

	res := make(chan bool)
	validationResults := make(chan []validationResult)
	go processResults(o, validationResults, res)

	c := cache.New()
	var wg sync.WaitGroup
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for fileBatch := range fileBatches {
				for _, filename := range fileBatch {
					f, err := os.Open(filename)
					if err != nil {
						fmt.Printf("failed opening %s\n", filename)
						validationResults <- []validationResult{{
							filename: filename,
							err:      err,
							skipped:  true,
						}}
						continue
					}

					res := ValidateStream(f, registries, k8sVersion, c, filter, ignoreMissingSchemas)
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
	close(validationResults)
	success := <-res
	if err = o.Flush(); err != nil {
		fmt.Fprint(os.Stderr, "failed flushing output")
	}

	if !success {
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
