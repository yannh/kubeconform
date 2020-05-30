package output

import (
	"github.com/yannh/kubeconform/pkg/validator"
	"log"
)

type TextOutput struct {
}

func NewTextOutput() Output {
	return &TextOutput{}
}

func (o *TextOutput) Write(filename string,err error, skipped bool) {
	if skipped {
		log.Printf("skipping resource\n")
		return
	}

	if err != nil {
		if _, ok := err.(validator.InvalidResourceError); ok {
			log.Printf("invalid resource: %s\n", err)
		} else {
			log.Printf("failed validating resource in file %s: %s\n", filename, err)
		}
		return
	}

	if !skipped{
		log.Printf("file %s is valid\n", filename)
	}
}

func (o *TextOutput) Flush() {
}

