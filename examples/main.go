package main

// WARNING: API of Kubeconform is still under development and not yet
// considered stable

import (
	"github.com/yannh/kubeconform/pkg/validator"
	"log"
	"os"
)

func main() {
	filepath := "../fixtures/valid.yaml"
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed opening %s: %s", filepath, err)
	}

	v := validator.New(nil, validator.Opts{Strict: true})
	for i, res := range v.Validate(filepath, f) { // A file might contain multiple resources
		// File starts with ---, the parser assumes a first empty resource
		if res.Status != validator.Valid && res.Status != validator.Empty {
			log.Fatalf("resource %d in file %s is not valid: %d, %s", i, filepath, res.Status, res.Err)
		}
	}
}
