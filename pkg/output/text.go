package output

import (
	"fmt"
	"github.com/yannh/kubeconform/pkg/validator"
)

type TextOutput struct {
	withSummary                bool
	nValid, nInvalid, nErrors, nSkipped int
}

func NewTextOutput(withSummary bool) Output {
	return &TextOutput{withSummary, 0,0,0, 0}
}

func (o *TextOutput) Write(filename string,err error, skipped bool) {
	if skipped {
		fmt.Printf("skipping resource\n")
		o.nSkipped++
		return
	}

	if err != nil {
		if _, ok := err.(validator.InvalidResourceError); ok {
			fmt.Printf("invalid resource: %s\n", err)
			o.nInvalid++
		} else {
			fmt.Printf("failed validating resource in file %s: %s\n", filename, err)
			o.nErrors++
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
		fmt.Printf("Run summary - Valid: %d, Invalid: %d, Errors: %d Skipped: %d\n", o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
	}
}

