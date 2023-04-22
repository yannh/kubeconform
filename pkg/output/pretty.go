package output

import (
	"fmt"
	"io"
	"sync"

	"github.com/yannh/kubeconform/pkg/validator"
)

type prettyo struct {
	sync.Mutex
	w                                   io.Writer
	withSummary                         bool
	isStdin                             bool
	verbose                             bool
	files                               map[string]bool
	nValid, nInvalid, nErrors, nSkipped int
}

// Text will output the results of the validation as a texto
func prettyOutput(w io.Writer, withSummary, isStdin, verbose bool) Output {
	return &prettyo{
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

func (o *prettyo) Write(result validator.Result) error {
	checkmark := "\u2714"
	multiplicationSign := "\u2716"
	reset := "\033[0m"
	cRed := "\033[31m"
	cGreen := "\033[32m"
	cYellow := "\033[33m"

	o.Lock()
	defer o.Unlock()

	var err error

	sig, _ := result.Resource.Signature()

	o.files[result.Resource.Path] = true
	switch result.Status {
	case validator.Valid:
		if o.verbose {
			fmt.Fprintf(o.w, "%s%s%s %s: %s%s %s is valid%s\n", cGreen, checkmark, reset, result.Resource.Path, cGreen, sig.Kind, sig.Name, reset)
		}
		o.nValid++
	case validator.Invalid:
		fmt.Fprintf(o.w, "%s%s%s %s: %s%s %s is invalid: %s%s\n", cRed, multiplicationSign, reset, result.Resource.Path, cRed, sig.Kind, sig.Name, result.Err.Error(), reset)

		o.nInvalid++
	case validator.Error:
		fmt.Fprintf(o.w, "%s%s%s %s: ", cRed, multiplicationSign, reset, result.Resource.Path)
		if sig.Kind != "" && sig.Name != "" {
			fmt.Fprintf(o.w, "%s%s failed validation: %s %s%s\n", cRed, sig.Kind, sig.Name, result.Err.Error(), reset)
		} else {
			fmt.Fprintf(o.w, "%sfailed validation: %s %s%s\n", cRed, sig.Name, result.Err.Error(), reset)
		}
		o.nErrors++
	case validator.Skipped:
		if o.verbose {
			fmt.Fprintf(o.w, "%s-%s %s: ", cYellow, reset, result.Resource.Path)
			if sig.Kind != "" && sig.Name != "" {
				fmt.Fprintf(o.w, "%s%s %s skipped%s\n", cYellow, sig.Kind, sig.Name, reset)
			} else if sig.Kind != "" {
				fmt.Fprintf(o.w, "%s%s skipped%s\n", cYellow, sig.Kind, reset)
			} else {
				fmt.Fprintf(o.w, "%sskipped%s\n", cYellow, reset)
			}
		}
		o.nSkipped++
	case validator.Empty: // sent to ensure we count the filename as parsed
	}

	return err
}

func (o *prettyo) Flush() error {
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
