package output

import (
	"fmt"
	"io"
	"sync"

	"github.com/logrusorgru/aurora"
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
	o.Lock()
	defer o.Unlock()

	var err error

	sig, _ := result.Resource.Signature()

	o.files[result.Resource.Path] = true
	switch result.Status {
	case validator.Valid:
		if o.verbose {
			fmt.Printf("%s %s: ", aurora.BrightGreen("\u2714"), result.Resource.Path)
			fmt.Println(aurora.Sprintf(aurora.BrightGreen("%s %s is valid"), sig.Kind, sig.Name))
		}
		o.nValid++
	case validator.Invalid:
		_, err = fmt.Fprintf(o.w, "%s - %s %s is invalid: %s\n", result.Resource.Path, sig.Kind, sig.Name, result.Err)
		o.nInvalid++
	case validator.Error:
		fmt.Printf("%s %s: ", aurora.BrightRed("\u2716"), result.Resource.Path)
		if sig.Kind != "" && sig.Name != "" {
			fmt.Println(aurora.Sprintf(aurora.BrightRed("%s failed validation: %s %s"), sig.Kind, sig.Name, result.Err))
		} else {
			fmt.Println(aurora.Sprintf(aurora.BrightRed("failed validation: %s %s"), sig.Name, result.Err))
		}
		o.nErrors++
	case validator.Skipped:
		if o.verbose {
			fmt.Printf("%s %s: ", aurora.Yellow("-"), result.Resource.Path)
			if sig.Kind != "" && sig.Name != "" {
				fmt.Println(aurora.Sprintf(aurora.Yellow("%s %s skipped"), sig.Kind, sig.Name))
			} else if sig.Kind != "" {
				fmt.Println(aurora.Sprintf(aurora.Yellow("%s skipped"), sig.Kind))
			} else {
				fmt.Println(aurora.Yellow("skipped"))
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
