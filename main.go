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
func validateFile(r io.Reader, regs []registry.Registry, k8sVersion string, c *cache.SchemaCache, skip func(signature resource.Signature) bool) []validationResult {
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

		if skip(sig) {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil, skipped: true})
			continue
		}

		var ok bool

		cacheKey := cache.Key(sig.Kind, sig.Version, k8sVersion)
		schema, ok := c.Get(cacheKey)
		if !ok {
			schema, err = downloadSchema(regs, sig.Kind, sig.Version, k8sVersion)
			if err != nil {
				validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: err, skipped: false})
				continue
			} else if schema == nil { // skip if no schema was found, but there was no error
				validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil, skipped: true})
				c.Set(cacheKey, nil)
				continue
			} else {
				c.Set(cacheKey, schema)
			}
		}

		if err = validator.Validate(rawResource, schema); err != nil {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: err})
		} else {
			validationResults = append(validationResults, validationResult{kind: sig.Kind, version: sig.Version, err: nil})
		}
	}

	return validationResults
}

type arrayParam []string

func (ap *arrayParam) String() string {
	s := ""
	for _, param := range *ap {
		if s == "" {
			s += param
		} else {
			s += " - " + param
		}
	}

	return s
}

func (ap *arrayParam) Set(value string) error {
	*ap = append(*ap, value)
	return nil
}

func getLogger(outputFormat string, printSummary, quiet bool) (output.Output, error) {
	switch {
	case outputFormat == "text":
		return output.NewTextOutput(printSummary, quiet), nil
	case outputFormat == "json":
		return output.NewJSONOutput(printSummary, quiet), nil
	default:
		return nil, fmt.Errorf("-output must be text or json")
	}
}

func skipKindsMap(skipKindsCSV string) map[string]bool {
	splitKinds := strings.Split(skipKindsCSV, ",")
	skipKinds := map[string]bool{}
	for _, kind := range splitKinds {
		skipKinds[kind] = true
	}
	return skipKinds
}

func realMain() int {
	var files, dirs, schemas arrayParam
	var skipKindsCSV, k8sVersion, outputFormat string
	var printSummary, strict, quiet bool
	var nWorkers int
	var err error

	flag.StringVar(&k8sVersion, "k8sversion", "1.18.0", "version of Kubernetes to test against")
	flag.Var(&files, "file", "file to validate (can be specified multiple times)")
	flag.Var(&dirs, "dir", "directory to validate (can be specified multiple times)")
	flag.Var(&schemas, "schema", "file containing an additional Schema (can be specified multiple times)")
	flag.BoolVar(&printSummary, "printsummary", false, "print a summary at the end")
	flag.IntVar(&nWorkers, "workers", 4, "number of routines to run in parallel")
	flag.StringVar(&skipKindsCSV, "skipKinds", "", "comma-separated list of kinds to ignore")
	flag.BoolVar(&strict, "strict", false, "disallow additional properties not in schema")
	flag.StringVar(&outputFormat, "output", "text", "output format - text, json")
	flag.BoolVar(&quiet, "quiet", false, "quiet output - only print invalid files, and errors")
	flag.Parse()

	var o output.Output
	if o, err = getLogger(outputFormat, printSummary, quiet); err != nil {
		fmt.Println(err)
		return 1
	}

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

	validationResults := make(chan []validationResult)
	var logWG sync.WaitGroup
	logWG.Add(1)

	success := true
	go func() {
		defer logWG.Done()
		for results := range validationResults {
			for _, result := range results {
				if result.err != nil {
					success = false
				}

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
	logWG.Wait()
	o.Flush()

	if !success {
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
