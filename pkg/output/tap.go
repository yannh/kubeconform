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
	index                               int
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
func (o *tapo) Write(res validator.Result) error {
	o.index++

	if o.index == 1 {
		fmt.Fprintf(o.w, "TAP version 13\n")
	}

	switch res.Status {
	case validator.Valid:
		sig, _ := res.Resource.Signature()
		fmt.Fprintf(o.w, "ok %d - %s (%s)\n", o.index, res.Resource.Path, sig.QualifiedName())

	case validator.Invalid:
		sig, _ := res.Resource.Signature()
		fmt.Fprintf(o.w, "not ok %d - %s (%s): %s\n", o.index, res.Resource.Path, sig.QualifiedName(), res.Err.Error())

	case validator.Empty:
		fmt.Fprintf(o.w, "ok %d - %s (empty)\n", o.index, res.Resource.Path)

	case validator.Error:
		fmt.Fprintf(o.w, "not ok %d - %s: %s\n", o.index, res.Resource.Path, res.Err.Error())

	case validator.Skipped:
		sig, _ := res.Resource.Signature()
		fmt.Fprintf(o.w, "ok %d - %s (%s) # skip\n", o.index, res.Resource.Path, sig.QualifiedName())
	}

	return nil
}

// Flush outputs the results as JSON
func (o *tapo) Flush() error {
	fmt.Fprintf(o.w, "1..%d\n", o.index)

	return nil
}
