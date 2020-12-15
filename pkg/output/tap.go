package output

import (
	"fmt"
	"io"

	"github.com/yannh/kubeconform/pkg/validator"
)

type tapo struct {
	w                                   io.Writer
	withSummary                         bool
	verbose                             bool
	results                             []validator.Result
	nValid, nInvalid, nErrors, nSkipped int
}

func tapOutput(w io.Writer, withSummary bool, isStdin, verbose bool) Output {
	return &tapo{
		w:           w,
		withSummary: withSummary,
		verbose:     verbose,
		results:     []validator.Result{},
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

// JSON.Write will only write when JSON.Flush has been called
func (o *tapo) Write(result validator.Result) error {
	o.results = append(o.results, result)
	return nil
}

// Flush outputs the results as JSON
func (o *tapo) Flush() error {
	var err error

	fmt.Fprintf(o.w, "1..%d\n", len(o.results))
	for i, res := range o.results {
		switch res.Status {
		case validator.Valid:
			sig, _ := res.Resource.Signature()
			fmt.Fprintf(o.w, "ok %d - %s (%s)\n", i, res.Resource.Path, sig.Kind)

		case validator.Invalid:
			sig, _ := res.Resource.Signature()
			fmt.Fprintf(o.w, "not ok %d - %s (%s): %s\n", i, res.Resource.Path, sig.Kind, res.Err.Error())

		case validator.Empty:
			fmt.Fprintf(o.w, "ok %d - %s (empty)\n", i, res.Resource.Path)

		case validator.Error:
			fmt.Fprintf(o.w, "not ok %d - %s: %s\n", i, res.Resource.Path, res.Err.Error())

		case validator.Skipped:
			fmt.Fprintf(o.w, "ok %d #skip - %s\n", i, res.Resource.Path)
		}
	}
	if err != nil {
		return err
	}

	return nil
}
