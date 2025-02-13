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
    switch outputFormat {
    case "json":
        return jsonOutput(w, printSummary, isStdin, verbose), nil
    case "junit":
        return junitOutput(w, printSummary, isStdin, verbose), nil
    case "pretty":
        return prettyOutput(w, printSummary, isStdin, verbose), nil
    case "tap":
        return tapOutput(w, printSummary, isStdin, verbose), nil
    case "text":
        return textOutput(w, printSummary, isStdin, verbose), nil
    default:
        return nil, fmt.Errorf("'outputFormat' must be 'json', 'junit', 'pretty', 'tap' or 'text'")
    }
}

// Mock writer for testing purposes
type mockWriter struct{}

func (m *mockWriter) Write(p []byte) (n int, err error) {
    return len(p), nil
}
