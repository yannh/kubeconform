package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yannh/kubeconform/pkg/validator"
)

type oresult struct {
	Filename string `json:"filename"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Status   string `json:"status"`
	Msg      string `json:"msg"`
}

type jsono struct {
	w                                   io.Writer
	withSummary                         bool
	verbose                             bool
	results                             []oresult
	nValid, nInvalid, nErrors, nSkipped int
}

// JSON will output the results of the validation as a JSON
func jsonOutput(w io.Writer, withSummary bool, isStdin, verbose bool) Output {
	return &jsono{
		w:           w,
		withSummary: withSummary,
		verbose:     verbose,
		results:     []oresult{},
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

// JSON.Write will only write when JSON.Flush has been called
func (o *jsono) Write(result validator.Result) error {
	msg, st := "", ""

	switch result.Status {
	case validator.Valid:
		st = "statusValid"
		o.nValid++
	case validator.Invalid:
		st = "statusInvalid"
		msg = result.Err.Error()
		o.nInvalid++
	case validator.Error:
		st = "statusError"
		msg = result.Err.Error()
		o.nErrors++
	case validator.Skipped:
		st = "statusSkipped"
		o.nSkipped++
	case validator.Empty:
	}

	if o.verbose || (result.Status != validator.Valid && result.Status != validator.Skipped && result.Status != validator.Empty) {
		sig, _ := result.Resource.Signature()
		o.results = append(o.results, oresult{Filename: result.Resource.Path, Kind: sig.Kind, Name: sig.Name, Version: sig.Version, Status: st, Msg: msg})
	}

	return nil
}

// Flush outputs the results as JSON
func (o *jsono) Flush() error {
	var err error
	var res []byte

	if o.withSummary {
		jsonObj := struct {
			Resources []oresult `json:"resources"`
			Summary   struct {
				Valid   int `json:"valid"`
				Invalid int `json:"invalid"`
				Errors  int `json:"errors"`
				Skipped int `json:"skipped"`
			} `json:"summary"`
		}{
			Resources: o.results,
			Summary: struct {
				Valid   int `json:"valid"`
				Invalid int `json:"invalid"`
				Errors  int `json:"errors"`
				Skipped int `json:"skipped"`
			}{
				Valid:   o.nValid,
				Invalid: o.nInvalid,
				Errors:  o.nErrors,
				Skipped: o.nSkipped,
			},
		}

		res, err = json.MarshalIndent(jsonObj, "", "  ")
	} else {
		jsonObj := struct {
			Resources []oresult `json:"resources"`
		}{
			Resources: o.results,
		}

		res, err = json.MarshalIndent(jsonObj, "", "  ")
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(o.w, "%s\n", res)

	return nil
}
