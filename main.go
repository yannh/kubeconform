package main

import (
	"flag"
	"fmt"
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
	err     error
	skipped bool
}

// filter returns true if the file should be skipped
// Returning an array, this Reader might container multiple resources
func validateFile(f io.Reader, regs []registry.Registry, k8sVersion string, skip func(signature resource.Signature)bool) []validationResult {
	rawResource, err := ioutil.ReadAll(f)
	if err != nil {
		return []validationResult{{err: fmt.Errorf("failed reading file: %s", err)}}
	}

	sig, err := resource.SignatureFromBytes(rawResource)
	if err != nil {
		return []validationResult{{err: fmt.Errorf("error while parsing: %s", err)}}
	}

	if skip(sig) {
		return []validationResult{{err: nil, skipped: true}}
	}

	var schema []byte
	for _, reg := range regs {
		downloadSchema := cache.WithCache(reg.DownloadSchema)
		schema, err = downloadSchema(sig.Kind, sig.Version, k8sVersion)
		if err == nil {
			break
		}
	}
	if err != nil {
		return []validationResult{{err: fmt.Errorf("error while downloading schema for resource: %s", err)}}
	}

	if err = validator.Validate(rawResource, schema); err != nil {
		return []validationResult{{err: err}}
	}

	return []validationResult{{err: nil}}
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
	var printSummary bool
	var nWorkers int

	flag.BoolVar(&printSummary, "printsummary", false, "print a summary at the end")
	flag.Var(&files, "file", "file to validate (can be specified multiple times)")
	flag.Var(&dirs, "dir", "directory to validate (can be specified multiple times)")
	flag.Var(&schemas, "schema", "file containing an additional Schema (can be specified multiple times)")
	flag.IntVar(&nWorkers, "workers", 4, "number of routines to run in parallel")
	flag.StringVar(&k8sVersion, "k8sversion", "1.18.0", "version of Kubernetes to test against")
	flag.StringVar(&skipKinds, "skipKinds", "", "comma-separated list of kinds to ignore")
	flag.StringVar(&outputFormat, "output", "text", "output format - text, json")
	flag.Parse()

	var o output.Output
	switch {
	case outputFormat == "text":
		o = output.NewTextOutput(printSummary)
	case outputFormat == "json":
		o = output.NewJSONOutput(printSummary)
	default:
		log.Fatalf("-output must be text or json")
	}

	filter := func(signature resource.Signature) bool {
		kinds := strings.Split(skipKinds, ",")
		for _, kind := range kinds {
			if signature.Kind == kind {
				return true
			}
		}
		return false
	}

	registries := []registry.Registry{}
	registries = append(registries, registry.NewKubernetesRegistry(false))
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

	var wg sync.WaitGroup
	for i:=0; i<nWorkers; i++ {
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

					res := validateFile(f, registries, k8sVersion, filter)
					f.Close()


					for _, resourceValidation := range res {
						o.Write(filename, resourceValidation.err, resourceValidation.skipped)
					}
				}
			}
		}()
	}

	wg.Wait()
	o.Flush()

	return 0
}

func main() {
	os.Exit(realMain())
}