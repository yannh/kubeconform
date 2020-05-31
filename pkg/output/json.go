package output

import (
	"encoding/json"
	"fmt"
)

type result struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
	Msg      string `json:"msg"`
}

type JSONOutput struct {
	withSummary bool
	results     []result
}

func NewJSONOutput(withSummary bool) Output {
	return &JSONOutput{
		withSummary: withSummary,
		results:     []result{},
	}
}

func (o *JSONOutput) Write(filename string, err error, skipped bool) {
	msg, st := "", ""

	s := status(err, skipped)
	switch {
	case s == VALID:
		st = "VALID"
	case s == INVALID:
		st = "INVALID"
		msg = err.Error()
	case s == ERROR:
		st = "ERROR"
		msg = err.Error()
	case s == SKIPPED:
		st = "SKIPPED"
	}

	o.results = append(o.results, result{Filename: filename, Status: st, Msg: msg})
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
		}

		for _, r := range o.results {
			switch {
			case r.Status == "VALID":
				jsonObj.Summary.Valid++
			case r.Status == "INVALID":
				jsonObj.Summary.Invalid++
			case r.Status == "ERROR":
				jsonObj.Summary.Errors++
			case r.Status == "SKIPPED":
				jsonObj.Summary.Skipped++
			}
		}

		res, err = json.MarshalIndent(jsonObj, "", "  ")
	} else {
		jsonObj := struct {
			Resources []result
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
