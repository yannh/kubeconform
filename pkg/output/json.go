package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type result struct {
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
	results                             []result
	nValid, nInvalid, nErrors, nSkipped int
}

// JSON will output the results of the validation as a JSON
func JSON(w io.Writer, withSummary bool, verbose bool) Output {
	return &jsono{
		w:           w,
		withSummary: withSummary,
		verbose:     verbose,
		results:     []result{},
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

// JSON.Write will only write when JSON.Flush has been called
func (o *jsono) Write(filename, kind, name, version string, err error, skipped bool) error {
	msg, st := "", ""

	s := status(err, skipped)

	switch s {
	case VALID:
		st = "VALID"
		o.nValid++
	case INVALID:
		st = "INVALID"
		msg = err.Error()
		o.nInvalid++
	case ERROR:
		st = "ERROR"
		msg = err.Error()
		o.nErrors++
	case SKIPPED:
		st = "SKIPPED"
		o.nSkipped++
	}

	if o.verbose || (s != VALID && s != SKIPPED) {
		o.results = append(o.results, result{Filename: filename, Kind: kind, Name: name, Version: version, Status: st, Msg: msg})
	}

	return nil
}

// Flush outputs the results as JSON
func (o *jsono) Flush() error {
	var err error
	var res []byte

	if o.withSummary {
		jsonObj := struct {
			Resources []result `json:"resources"`
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
			Resources []result `json:"resources"`
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
