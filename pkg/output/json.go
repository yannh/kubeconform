package output

import (
	"encoding/json"
	"fmt"
)

type result struct {
	Filename string  `json:"filename"`
	Status string `json:"status"`
	Msg string `json:"msg"`
}

type JSONOutput struct {
	withSummary bool
	results []result
}

func NewJSONOutput(withSummary bool) Output{
	return &JSONOutput{
		withSummary: withSummary,
		results: []result{},
	}
}

func (o *JSONOutput) Write(filename string,err error, skipped bool) {
	status := "VALID"
	msg := ""
	if err != nil {
		status = "INVALID"
		msg = err.Error()
	}
	if skipped {
		status = "SKIPPED"
	}

	o.results = append(o.results, result{Filename: filename, Status: status, Msg: msg})
}

func (o *JSONOutput) Flush() {
	var err error
	var res []byte

	if o.withSummary {
		jsonObj := struct {
			Resources []result `json:"resources"`
			Summary struct {
				Valid int `json:"valid"`
				Invalid int `json:"invalid"`
				Skipped int `json:"skipped"`
			} `json:"summary"`
		} {
			Resources: o.results,
		}

		for _, r := range o.results {
			switch {
			case r.Status == "VALID":
				jsonObj.Summary.Valid++
			case r.Status == "INVALID":
				jsonObj.Summary.Invalid++
			case r.Status == "SKIPPED":
				jsonObj.Summary.Skipped++
			}
		}

		res, err = json.MarshalIndent(jsonObj,"", "  ")
	} else {
		jsonObj := struct {
			Resources []result
		} {
			Resources: o.results,
		}

		res, err = json.MarshalIndent(jsonObj,"", "  ")
	}

	if err != nil {
		fmt.Printf("error print results: %s", err)
		return
	}
	fmt.Printf("%s\n", res)
}
