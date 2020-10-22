package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/xeipuuv/gojsonschema"

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

type arrayParam []string

func (ap *arrayParam) String() string {
	return strings.Join(*ap, " - ")
}

func (ap *arrayParam) Set(value string) error {
	*ap = append(*ap, value)
	return nil
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
	var schemaLocationsParam arrayParam
	var skipKindsCSV, k8sVersion, outputFormat string
	var summary, strict, verbose, ignoreMissingSchemas, help bool
	var nWorkers int
	var err error
	var files []string

	flag.StringVar(&k8sVersion, "kubernetes-version", "1.18.0", "version of Kubernetes to validate against")
	flag.Var(&schemaLocationsParam, "schema-location", "override schemas location search path (can be specified multiple times)")
	flag.BoolVar(&ignoreMissingSchemas, "ignore-missing-schemas", false, "skip files with missing schemas instead of failing")
	flag.BoolVar(&summary, "summary", false, "print a summary at the end")
	flag.IntVar(&nWorkers, "n", 4, "number of routines to run in parallel")
	flag.StringVar(&skipKindsCSV, "skip", "", "comma-separated list of kinds to ignore")
	flag.BoolVar(&strict, "strict", false, "disallow additional properties not in schema")
	flag.StringVar(&outputFormat, "output", "text", "output format - text, json")
	flag.BoolVar(&verbose, "verbose", false, "print results for all resources")
	flag.BoolVar(&help, "h", false, "show help information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [FILE OR FOLDER]...\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if help {
		flag.Usage()
		return 1
	}

	skipKinds := skipKindsMap(skipKindsCSV)

	for _, file := range flag.Args() {
		files = append(files, file)
	}

	filter := func(signature resource.Signature) bool {
		isSkipKind, ok := skipKinds[signature.Kind]
		return ok && isSkipKind
	}

	if len(schemaLocationsParam) == 0 {
		schemaLocationsParam = append(schemaLocationsParam, "https://kubernetesjsonschema.dev") // if not specified, default behaviour is to use kubernetesjson-schema.dev as registry
	}

	registries := []registry.Registry{}
	for _, schemaLocation := range schemaLocationsParam {
		if !strings.HasSuffix(schemaLocation, "json") { // If we dont specify a full templated path, we assume the paths of kubernetesjsonschema.dev
			schemaLocation += "/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json"
		}

		if strings.HasPrefix(schemaLocation, "http") {
			registries = append(registries, registry.NewHTTPRegistry(schemaLocation, strict))
		} else {
			registries = append(registries, registry.NewLocalRegistry(schemaLocation, strict))
		}
	}

	validationResults := make(chan []validationResult)

	fileBatches := make(chan []string)
	go func() {
		getFiles(files, fileBatches, validationResults)
		close(fileBatches)
	}()

	var o output.Output
	if o, err = output.New(outputFormat, summary, verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	res := make(chan bool)
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
	o.Flush()

	if !success {
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
