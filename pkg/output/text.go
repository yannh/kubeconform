package output

import (
	"fmt"
	"io"
	"sync"

	"github.com/yannh/kubeconform/pkg/validator"
)

type texto struct {
	sync.Mutex
	w                                   io.Writer
	withSummary                         bool
	isStdin                             bool
	verbose                             bool
	files                               map[string]bool
	nValid, nInvalid, nErrors, nSkipped int
}

// Text will output the results of the validation as a texto
func textOutput(w io.Writer, withSummary, isStdin, verbose bool) Output {
	return &texto{
		w:           w,
		withSummary: withSummary,
		isStdin:     isStdin,
		verbose:     verbose,
		files:       map[string]bool{},
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

func (o *texto) Write(result validator.Result) error {
	o.Lock()
	defer o.Unlock()

	var err error

	sig, _ := result.Resource.Signature()

	o.files[result.Resource.Path] = true
	switch result.Status {
	case validator.Valid:
		if o.verbose {
			_, err = fmt.Fprintf(o.w, "%s - %s %s is valid\n", result.Resource.Path, sig.Kind, sig.Name)
		}
		o.nValid++
	case validator.Invalid:
		_, err = fmt.Fprintf(o.w, "%s - %s %s is invalid: %s\n", result.Resource.Path, sig.Kind, sig.Name, result.Err)
		o.nInvalid++
	case validator.Error:
		if sig.Kind != "" && sig.Name != "" {
			_, err = fmt.Fprintf(o.w, "%s - %s %s failed validation: %s\n", result.Resource.Path, sig.Kind, sig.Name, result.Err)
		} else {
			_, err = fmt.Fprintf(o.w, "%s - failed validation: %s\n", result.Resource.Path, result.Err)
		}
		o.nErrors++
	case validator.Skipped:
		if o.verbose {
			_, err = fmt.Fprintf(o.w, "%s - %s %s skipped\n", result.Resource.Path, sig.Name, sig.Kind)
		}
		o.nSkipped++
	case validator.Empty: // sent to ensure we count the filename as parsed
	}

	return err
}

func (o *texto) Flush() error {
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
		if o.isStdin {
			_, err = fmt.Fprintf(o.w, "Summary: %d resource%s found parsing stdin - Valid: %d, Invalid: %d, Errors: %d, Skipped: %d\n", nResources, resourcesPlural, o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
		} else {
			_, err = fmt.Fprintf(o.w, "Summary: %d resource%s found in %d file%s - Valid: %d, Invalid: %d, Errors: %d, Skipped: %d\n", nResources, resourcesPlural, nFiles, filesPlural, o.nValid, o.nInvalid, o.nErrors, o.nSkipped)
		}
	}

	return err
}
