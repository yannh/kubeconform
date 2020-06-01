package output

import (
	"fmt"
	"io"
	"sync"
)

type text struct {
	sync.Mutex
	w                                   io.Writer
	withSummary                         bool
	verbose                             bool
	nValid, nInvalid, nErrors, nSkipped int
}

func Text(w io.Writer, withSummary, verbose bool) Output {
	return &text{
		w:           w,
		withSummary: withSummary,
		verbose:     verbose,
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

func (o *text) Write(filename, kind, version string, err error, skipped bool) {
	o.Lock()
	defer o.Unlock()

	switch status(err, skipped) {
	case VALID:
		if !o.verbose {
			fmt.Fprintf(o.w, "%s - %s is valid\n", filename, kind)
		}
		o.nValid++
	case INVALID:
		fmt.Fprintf(o.w, "%s - %s is invalid: %s\n", filename, kind, err)
		o.nInvalid++
	case ERROR:
		fmt.Fprintf(o.w, "%s - %s failed validation: %s\n", filename, kind, err)
		o.nErrors++
	case SKIPPED:
		if o.verbose {
			fmt.Fprintf(o.w, "%s - %s skipped\n", filename, kind)
		}
		o.nSkipped++
	}
}

func (o *text) Flush() {
	if o.withSummary {
		fmt.Fprintf(o.w, "Run summary - Valid: %d, Invalid: %d, Errors: %d Skipped: %d\n", o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
	}
}
