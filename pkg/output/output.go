package output

import (
	"fmt"
	"github.com/yannh/kubeconform/pkg/validator"
	"os"
)

const (
	_ = iota
	statusValid
	statusInvalid
	statusError
	statusSkipped
	statusEmpty
)

type Output interface {
	Write(filename, kind, name, version string, err error, skipped bool) error
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

func status(kind, name string, err error, skipped bool) int {
	if name == "" && kind == "" && err == nil && skipped == false {
		return statusEmpty
	}

	if skipped {
		return statusSkipped
	}

	if err != nil {
		if _, ok := err.(validator.InvalidResourceError); ok {
			return statusInvalid
		}
		return statusError
	}

	return statusValid
}
