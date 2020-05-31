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

// filter returns true if the file should be skipped
// Returning an array, this Reader might container multiple resources
func validateFile(f io.Reader, regs []registry.Registry, k8sVersion string, c *cache.SchemaCache, skip func(signature resource.Signature) bool) []validationResult {
	file, err := ioutil.ReadAll(f)
	if err != nil {
		return []validationResult{{err: fmt.Errorf("failed reading file: %s", err)}}
	}

	validationResults := []validationResult{}
	rawResources := bytes.Split(file, []byte("---\n"))

RESOURCES:
	for _, rawResource := range rawResources {
		sig, err := resource.SignatureFromBytes(rawResource)
		if err != nil {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: fmt.Errorf("error while parsing: %s", err)})
			continue
		}

		if skip(sig) {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil, skipped: true})
			continue
		}

		var schema *gojsonschema.Schema
		var schemaBytes []byte
		var ok bool

		cacheKey := cache.Key(sig.Kind, sig.Version, k8sVersion)
		schema, ok = c.Get(cacheKey)
		if !ok {
			for _, reg := range regs {
				schemaBytes, err = reg.DownloadSchema(sig.Kind, sig.Version, k8sVersion)
				if err == nil {
					schema, err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaBytes))
					if err != nil {
						// Downloaded a schema but failed to parse it
						validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: err, skipped: false})
						continue RESOURCES
					}

					// success
					break
				}

				// If we get a 404, we keep trying, but we exit if we get a real failure
				if er, retryable := err.(registry.Retryable); !(retryable && !er.IsRetryable()) {
					validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: fmt.Errorf("error while downloading schema for resource: %s", err)})
				}
			}
		}

		// Cache both found & not found
		c.Set(cacheKey, schema)

		if err != nil { // Not found
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil, skipped: true}) // skip if no schema found
		}

		if err = validator.Validate(rawResource, schema); err != nil {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: err})
		}

		validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil})
	}

	return validationResults
}

type arrayFiles []string

func (i *arrayFiles) String() string {
	return "my string representation"
}

func (i *arrayFiles) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func realMain() int {
	var files, dirs, schemas arrayFiles
	var skipKinds, k8sVersion, outputFormat string
	var printSummary, strict bool
	var nWorkers int
	var quiet bool

	flag.StringVar(&k8sVersion, "k8sversion", "1.18.0", "version of Kubernetes to test against")
	flag.Var(&files, "file", "file to validate (can be specified multiple times)")
	flag.Var(&dirs, "dir", "directory to validate (can be specified multiple times)")
	flag.Var(&schemas, "schema", "file containing an additional Schema (can be specified multiple times)")
	flag.BoolVar(&printSummary, "printsummary", false, "print a summary at the end")
	flag.IntVar(&nWorkers, "workers", 4, "number of routines to run in parallel")
	flag.StringVar(&skipKinds, "skipKinds", "", "comma-separated list of kinds to ignore")
	flag.BoolVar(&strict, "strict", false, "disallow additional properties not in schema")
	flag.StringVar(&outputFormat, "output", "text", "output format - text, json")
	flag.BoolVar(&quiet, "quiet", false, "quiet output - only print invalid files, and errors")
	flag.Parse()

	var o output.Output
	switch {
	case outputFormat == "text":
		o = output.NewTextOutput(printSummary, quiet)
	case outputFormat == "json":
		o = output.NewJSONOutput(printSummary, quiet)
	default:
		log.Fatalf("-output must be text or json")
	}

	splitKinds := strings.Split(skipKinds, ",")
	kinds := map[string]bool{}
	for _, kind := range splitKinds {
		kinds[kind] = true
	}
	filter := func(signature resource.Signature) bool {
		isSkipKind, ok := kinds[signature.Kind]
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

	validationResults := make(chan []validationResult)
	var logWG sync.WaitGroup
	logWG.Add(1)
	go func() {
		defer logWG.Done()
		for results := range validationResults {
			for _, result := range results {
				o.Write(result.filename, result.kind, result.version, result.err, result.skipped)
			}
		}
	}()

	c := cache.NewSchemaCache()
	var wg sync.WaitGroup
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for fileBatch := range fileBatches {
				for _, filename := range fileBatch {
					f, err := os.Open(filename)
					if err != nil {
						log.Printf("failed opening %s\n", filename)
						continue
					}

					res := validateFile(f, registries, k8sVersion, c, filter)
					f.Close()

					for i, _ := range res {
						res[i].filename = filename
					}
					validationResults <- res
				}
			}
		}()
	}

	wg.Wait()
	o.Flush()
	close(validationResults)
	logWG.Wait()

	return 0
}

func main() {
	os.Exit(realMain())
}
