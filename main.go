package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

func validateFile(f io.Reader, regs []*registry.KubernetesRegistry, k8sVersion string) error {
	rawResource, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed reading file: %s", err)
	}

	sig, err := resource.SignatureFromBytes(rawResource)
	if err != nil {
		return fmt.Errorf("error while parsing: %s", err)
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
		return fmt.Errorf("failed downloading schema for resource")
	}

	if err = validator.Validate(rawResource, schema); err != nil {
		return err
	}

	return nil
}

func realMain() int {
	const k8sVersion = "1.18.0"
	filename := "fixtures/valid_1.yaml"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed opening %s", filename)
		return 1
	}
	defer f.Close()

	r := registry.NewKubernetesRegistry()
	if err = validateFile(f, []*registry.KubernetesRegistry{r}, k8sVersion); err != nil {
		if _, ok := err.(validator.InvalidResourceError); ok {
			log.Printf("invalid resource: %s", err)
			return 1
		}
		log.Printf("failed validating resource: %s", err)
		return 1
	}

	log.Printf("resource is valid: %s", filename)
	return 0
}

func main() {
	os.Exit(realMain())
}