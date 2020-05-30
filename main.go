package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

func realMain() int {
	const k8sVersion = "1.18.0"
	filename := "fixtures/valid_1.yaml"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed opening %s", filename)
		return 1
	}
	defer f.Close()

	rawResource, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("failed reading file %s", filename)
		return 1
	}

	sig, err := resource.SignatureFromBytes(rawResource)
	if err != nil {
		log.Printf("failed parsing %s", filename)
		return 1
	}

	r := registry.NewKubernetesRegistry()
	schema, err := r.DownloadSchema(sig.Kind, sig.Version, k8sVersion)

	if err != nil {
		log.Printf("error downloading Schema: %s", err)
		return 1
	}

	err = validator.Validate(rawResource, schema)
	if err != nil {
		log.Printf("failed validating: %s", err)
		return 1
	}


	log.Printf("resource is valid!: %s", schema)

	return 0
}

func main() {
	os.Exit(realMain())
}