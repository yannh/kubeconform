package output

import (
	"fmt"
	"github.com/yannh/kubeconform/pkg/validator"
	"os"
)

type Output interface {
	Write(validator.Result) error
	Flush() error
}

func New(outputFormat string, printSummary, isStdin, verbose bool) (Output, error) {
	w := os.Stdout

	switch {
	case outputFormat == "text":
		return textOutput(w, printSummary, isStdin, verbose), nil
	case outputFormat == "json":
		return jsonOutput(w, printSummary, isStdin, verbose), nil
	default:
		return nil, fmt.Errorf("`outputFormat` must be 'text' or 'json'")
	}
}
