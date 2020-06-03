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

// Text will output the results of the validation as a text
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

func (o *text) Write(filename, kind, version string, reserr error, skipped bool) error {
	o.Lock()
	defer o.Unlock()

	var err error

	switch status(reserr, skipped) {
	case VALID:
		if o.verbose {
			_, err = fmt.Fprintf(o.w, "%s - %s is valid\n", filename, kind)
		}
		o.nValid++
	case INVALID:
		_, err = fmt.Fprintf(o.w, "%s - %s is invalid: %s\n", filename, kind, reserr)
		o.nInvalid++
	case ERROR:
		_, err = fmt.Fprintf(o.w, "%s - %s failed validation: %s\n", filename, kind, reserr)
		o.nErrors++
	case SKIPPED:
		if o.verbose {
			_, err = fmt.Fprintf(o.w, "%s - %s skipped\n", filename, kind)
		}
		o.nSkipped++
	}

	return err
}

func (o *text) Flush() error {
	var err error
	if o.withSummary {
		_, err = fmt.Fprintf(o.w, "Run summary - Valid: %d, Invalid: %d, Errors: %d Skipped: %d\n", o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
	}

	return err
}
