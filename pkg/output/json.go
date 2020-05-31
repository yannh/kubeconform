package output

import (
	"encoding/json"
	"fmt"
	"sync"
)

type result struct {
	Filename string `json:"filename"`
	Kind     string `json:"kind"`
	Version  string `json:"version"`
	Status   string `json:"status"`
	Msg      string `json:"msg"`
}

type JSONOutput struct {
	sync.Mutex
	withSummary                         bool
	quiet                               bool
	results                             []result
	nValid, nInvalid, nErrors, nSkipped int
}

func NewJSONOutput(withSummary bool, quiet bool) Output {
	return &JSONOutput{
		withSummary: withSummary,
		quiet:       quiet,
		results:     []result{},
		nValid:      0,
		nInvalid:    0,
		nErrors:     0,
		nSkipped:    0,
	}
}

func (o *JSONOutput) Write(filename, kind, version string, err error, skipped bool) {
	o.Lock()
	defer o.Unlock()
	msg, st := "", ""

	s := status(err, skipped)
	switch {
	case s == VALID:
		st = "VALID"
		o.nValid++
	case s == INVALID:
		st = "INVALID"
		msg = err.Error()
		o.nInvalid++
	case s == ERROR:
		st = "ERROR"
		msg = err.Error()
		o.nErrors++
	case s == SKIPPED:
		st = "SKIPPED"
		o.nSkipped++
	}

	if !o.quiet || (s != VALID && s != SKIPPED) {
		o.results = append(o.results, result{Filename: filename, Kind: kind, Version: version, Status: st, Msg: msg})
	}
}

func (o *JSONOutput) Flush() {
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
		fmt.Printf("error print results: %s", err)
		return
	}
	fmt.Printf("%s\n", res)
}
