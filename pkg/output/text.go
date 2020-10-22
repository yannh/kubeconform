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
	files                               map[string]bool
	nValid, nInvalid, nErrors, nSkipped int
}

// Text will output the results of the validation as a text
func Text(w io.Writer, withSummary, verbose bool) Output {
	return &text{
		w:           w,
		withSummary: withSummary,
		verbose:     verbose,
		files:       map[string]bool{},
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

func (o *text) Write(filename, kind, name, version string, reserr error, skipped bool) error {
	o.Lock()
	defer o.Unlock()

	var err error

	o.files[filename] = true
	switch status(kind, name, reserr, skipped) {
	case statusValid:
		if o.verbose {
			_, err = fmt.Fprintf(o.w, "%s - %s %s is valid\n", filename, kind, name)
		}
		o.nValid++
	case statusInvalid:
		_, err = fmt.Fprintf(o.w, "%s - %s %s is invalid: %s\n", filename, kind, name, reserr)
		o.nInvalid++
	case statusError:
		if kind != "" && name != "" {
			_, err = fmt.Fprintf(o.w, "%s - %s %s failed validation: %s\n", filename, kind, name, reserr)
		} else {
			_, err = fmt.Fprintf(o.w, "%s - failed validation: %s\n", filename, reserr)
		}
		o.nErrors++
	case statusSkipped:
		if o.verbose {
			_, err = fmt.Fprintf(o.w, "%s - %s %s skipped\n", filename, name, kind)
		}
		o.nSkipped++
	case statusEmpty: // sent to ensure we count the filename as parsed
	}

	return err
}

func (o *text) Flush() error {
	var err error
	if o.withSummary {
		nFiles := len(o.files)
		nResources := o.nValid + o.nInvalid + o.nErrors + o.nSkipped
		resourcesPlural := ""
		if nResources > 1 {
			resourcesPlural = "s"
		}
		filesPlural := ""
		if nFiles > 1 {
			filesPlural = "s"
		}
		_, err = fmt.Fprintf(o.w, "Summary: %d resource%s found in %d file%s - Valid: %d, Invalid: %d, Errors: %d Skipped: %d\n", nResources, resourcesPlural, nFiles, filesPlural, o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
	}

	return err
}
