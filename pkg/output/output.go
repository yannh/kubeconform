package output

import (
	"fmt"
	"os"
)

// TODO comment
const (
	_ = iota
	VALID
	INVALID
	ERROR
	SKIPPED
	EMPTY
)

// Output TODO
type Output interface {
	Write(filename, kind, name, version string, err error, skipped bool) error
	Flush() error
}

// New TODO
func New(outputFormat string, printSummary, verbose bool) (Output, error) {
	w := os.Stdout

	switch {
	case outputFormat == "text":
		return Text(w, printSummary, verbose), nil
	case outputFormat == "json":
		return JSON(w, printSummary, verbose), nil
	default:
		return nil, fmt.Errorf("-output must be text or json")
	}
}

func status(kind, name string, err error, skipped bool) int {
	if name == "" && kind == "" && err == nil && skipped == false {
		return EMPTY
	}

	if skipped {
		return SKIPPED
	}

	if err != nil {
		// if _, ok := err.(validator.InvalidResourceError); ok {
		// return INVALID
		// }
		return ERROR
	}

	return VALID
}
