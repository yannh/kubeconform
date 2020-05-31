package output

import (
	"fmt"
)

type TextOutput struct {
	withSummary                         bool
	nValid, nInvalid, nErrors, nSkipped int
}

func NewTextOutput(withSummary bool) Output {
	return &TextOutput{withSummary, 0, 0, 0, 0}
}

func (o *TextOutput) Write(filename string, err error, skipped bool) {
	s := status(err, skipped)
	switch {
	case s == VALID:
		fmt.Printf("file %s is valid\n", filename)
		o.nValid++
	case s == INVALID:
		fmt.Printf("invalid resource: %s\n", err)
		o.nInvalid++
	case s == ERROR:
		fmt.Printf("failed validating resource in file %s: %s\n", filename, err)
		o.nErrors++
	case s == SKIPPED:
		fmt.Printf("skipping resource\n")
		o.nSkipped++
	}
}

func (o *TextOutput) Flush() {
	if o.withSummary {
		fmt.Printf("Run summary - Valid: %d, Invalid: %d, Errors: %d Skipped: %d\n", o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
	}
}
