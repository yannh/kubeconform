package main

import (
	"log"
	"os"

	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
)

func realMain() int {
	filename := "fixtures/valid_1.yaml"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed opening %s", filename)
		return 1
	}
	defer f.Close()

	res, err := resource.Read(f)
	if err != nil {
		log.Printf("failed parsing %s", filename)
		return 1
	}

	r := registry.NewKubernetesRegistry()
	schema, err := r.DownloadSchema(res.Kind, res.Version, "1.18.0")

	if err != nil {
		log.Printf("error downloading Schema: %s")
		return 1
	}

	log.Printf("downloaded schema successfully: %s", schema)

	return 0
}

func main() {
	os.Exit(realMain())
}