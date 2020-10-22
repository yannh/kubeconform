package output

import (
	"fmt"
	"os"

	"github.com/yannh/kubeconform/pkg/validator"
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

func New(outputFormat string, printSummary, verbose bool) (Output, error) {
	w := os.Stdout

	switch {
	case outputFormat == "text":
		return Text(w, printSummary, verbose), nil
	case outputFormat == "json":
		return JSON(w, printSummary, verbose), nil
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
