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

func (o *TextOutput) Write(filename, kind, version string, err error, skipped bool) {
	o.Lock()
	defer o.Unlock()

	switch status(err, skipped) {
	case VALID:
		if !o.quiet {
			fmt.Printf("%s - %s is valid\n", filename, kind)
		}
		o.nValid++
	case INVALID:
		fmt.Printf("%s - %s is invalid: %s\n", filename, kind, err)
		o.nInvalid++
	case ERROR:
		fmt.Printf("%s - %s failed validation: %s\n", filename, kind, err)
		o.nErrors++
	case SKIPPED:
		if !o.quiet {
			fmt.Printf("%s - %s skipped\n", filename, kind)
		}
		o.nSkipped++
	}
}

func (o *TextOutput) Flush() {
	if o.withSummary {
		fmt.Printf("Run summary - Valid: %d, Invalid: %d, Errors: %d Skipped: %d\n", o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
	}
}
