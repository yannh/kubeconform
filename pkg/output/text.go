package output

import (
	"fmt"
	"sync"
)

type TextOutput struct {
	sync.Mutex
	withSummary                         bool
	quiet                               bool
	nValid, nInvalid, nErrors, nSkipped int
}

func NewTextOutput(withSummary, quiet bool) Output {
	return &TextOutput{
		withSummary: withSummary,
		quiet:       quiet,
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

func (o *TextOutput) Write(filename string, err error, skipped bool) {
	o.Lock()
	defer o.Unlock()

	s := status(err, skipped)
	switch {
	case s == VALID:
		if !o.quiet {
			fmt.Printf("file %s is valid\n", filename)
		}
		o.nValid++
	case s == INVALID:
		fmt.Printf("invalid resource: %s\n", err)
		o.nInvalid++
	case s == ERROR:
		fmt.Printf("failed validating resource in file %s: %s\n", filename, err)
		o.nErrors++
	case s == SKIPPED:
		if !o.quiet {
			fmt.Printf("skipping resource\n")
		}
		o.nSkipped++
	}
}

func (o *TextOutput) Flush() {
	if o.withSummary {
		fmt.Printf("Run summary - Valid: %d, Invalid: %d, Errors: %d Skipped: %d\n", o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
	}
}
