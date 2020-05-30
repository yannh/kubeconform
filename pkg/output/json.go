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
	results []result
}

func NewJSONOutput() Output{
	return &JSONOutput{
		results: []result{},
	}
}

func (o *JSONOutput) Write(filename string,err error, skipped bool) {
	status := "VALID"
	msg := ""
	if err != nil {
		status = "ERROR"
		msg = err.Error()
	}
	if skipped {
		status = "SKIPPED"
	}

	o.results = append(o.results, result{Filename: filename, Status: status, Msg: msg})
}

func (o *JSONOutput) Flush() {
	res, err := json.MarshalIndent(o.results,"", "  ")
	if err != nil {
		fmt.Printf("error print results: %s", err)
		return
	}
	fmt.Printf("%s\n", res)
}
