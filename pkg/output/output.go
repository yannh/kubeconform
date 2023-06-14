package output

import (
	"fmt"
	"io"

	"github.com/yannh/kubeconform/pkg/validator"
)

type Output interface {
	Write(validator.Result) error
	Flush() error
}

func New(w io.Writer, outputFormat string, printSummary, isStdin, verbose bool) (Output, error) {
	switch {
	case outputFormat == "json":
		return jsonOutput(w, printSummary, isStdin, verbose), nil
	case outputFormat == "junit":
		return junitOutput(w, printSummary, isStdin, verbose), nil
	case outputFormat == "pretty":
		return prettyOutput(w, printSummary, isStdin, verbose), nil
	case outputFormat == "tap":
		return tapOutput(w, printSummary, isStdin, verbose), nil
	case outputFormat == "text":
		return textOutput(w, printSummary, isStdin, verbose), nil
	default:
		return nil, fmt.Errorf("'outputFormat' must be 'json', 'junit', 'pretty', 'tap' or 'text'")
	}
}
