package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
func validateFile(f io.Reader, regs []*registry.KubernetesRegistry, k8sVersion string, skip func(signature resource.Signature)bool) []validationResult {
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
	var files, dirs arrayFiles
	var skipKinds, k8sVersion string
	flag.Var(&files, "file", "file to validate (can be specified multiple times)")
	flag.Var(&dirs, "dir", "directory to validate (can be specified multiple times)")
	flag.StringVar(&k8sVersion, "k8sversion", "1.18.0", "version of Kubernetes to test against")
	flag.StringVar(&skipKinds, "skipKinds", "", "comma-separated list of kinds to ignore")
	flag.Parse()

	filter := func(signature resource.Signature) bool {
		kinds := strings.Split(skipKinds, ",")
		for _, kind := range kinds {
			if signature.Kind == kind {
				return true
			}
		}
		return false
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

	r := registry.NewKubernetesRegistry(false)

	for fileBatch := range fileBatches {
		for _, filename := range fileBatch {
			f, err := os.Open(filename)
			if err != nil {
				log.Printf("failed opening %s\n", filename)
				continue
			}

			res := validateFile(f, []*registry.KubernetesRegistry{r}, k8sVersion, filter)
			f.Close()
			for _, resourceValidation := range res {
				if resourceValidation.skipped {
					log.Printf("skipping resource\n")
					continue
				}

				if resourceValidation.err != nil {
					if _, ok := resourceValidation.err.(validator.InvalidResourceError); ok {
						log.Printf("invalid resource: %s\n", resourceValidation.err)
					} else {
						log.Printf("failed validating resource in file %s: %s\n", filename, resourceValidation.err)
					}
					continue
				}

				if resourceValidation.err == nil && !resourceValidation.skipped{
					log.Printf("file %s is valid\n", filename)
				}
			}
		}
	}

	return 0
}

func main() {
	os.Exit(realMain())
}