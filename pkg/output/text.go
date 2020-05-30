package output

import (
	"fmt"
	"github.com/yannh/kubeconform/pkg/validator"
)

type TextOutput struct {
	withSummary                bool
	nValid, nInvalid, nSkipped int
}

func NewTextOutput(withSummary bool) Output {
	return &TextOutput{withSummary, 0,0,0}
}

func (o *TextOutput) Write(filename string,err error, skipped bool) {
	if skipped {
		fmt.Printf("skipping resource\n")
		o.nSkipped++
		return
	}

	if err != nil {
		o.nInvalid++
		if _, ok := err.(validator.InvalidResourceError); ok {
			fmt.Printf("invalid resource: %s\n", err)
		} else {
			fmt.Printf("failed validating resource in file %s: %s\n", filename, err)
		}
		return
	}

	if !skipped{
		fmt.Printf("file %s is valid\n", filename)
		o.nValid++
	}
}

func (o *TextOutput) Flush() {
	if o.withSummary {
		fmt.Printf("Run summary - Valid: %d, Invalid: %d, Skipped: %d\n", o.nValid, o.nInvalid, o.nSkipped)
	}
}

