package sarif

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

type Report struct {
	// The URI of the JSON schema corresponding to the version.
	Schema string `json:"$schema"`

	// The SARIF format version of this log file.
	Version string `json:"version"`

	// The set of runs contained in this log file.
	Runs []*Run `json:"runs"`

	// References to external property files that should be inlined with the content of a root log file.
	InlineExternalProperties []*ExternalProperties `json:"inlineExternalProperties,omitempty"`

	// Key/value pairs that provide additional information about the report.
	Properties PropertyBag `json:"properties,omitempty"`
}

func NewReport() *Report {
	return &Report{
		Schema:                   "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json",
		Version:                  "2.1.0",
		Runs:                     []*Run{},
		InlineExternalProperties: []*ExternalProperties{},
		Properties:               *NewPropertyBag(),
	}
}

// AddRun adds a run to the report
func (r *Report) AddRun(run *Run) *Report {
	r.Runs = append(r.Runs, run)
	return r
}

// Validate validates the report against the SARIF schema
func (r *Report) Validate() error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewGoLoader(r)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	var errors []string

	if !result.Valid() {
		for _, desc := range result.Errors() {
			errors = append(errors, fmt.Sprintf("%s\n", desc.String()))
		}
		return fmt.Errorf("validation failed: %v", errors)
	}

	return nil
}

// NewRunWithInformationURI creates a new Run and returns a pointer to it
func NewRunWithInformationURI(toolName, informationURI string) *Run {
	run := NewRun()
	run.Tool = NewTool()
	run.Tool.Driver = NewToolComponent()
	run.Tool.Driver.Name = &toolName
	run.Tool.Driver.InformationURI = &informationURI
	return run
}

// Open loads a Report from a file
func Open(filename string) (*Report, error) {
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return nil, errors.New("the provided file path doesn't have a file")
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("the provided filepath could not be opened. %w", err)
	}
	return FromBytes(content)
}

// FromString loads a Report from string content
func FromString(content string) (*Report, error) {
	return FromBytes([]byte(content))
}

// FromBytes loads a Report from a byte array
func FromBytes(content []byte) (*Report, error) {
	var report Report
	if err := json.Unmarshal(content, &report); err != nil {
		return nil, err
	}
	return &report, nil
}

// WriteFile will write the report to a file using a pretty formatter
func (sarif *Report) WriteFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	return sarif.PrettyWrite(file)
}

// Write writes the JSON as a string with no formatting
func (sarif *Report) Write(w io.Writer) error {
	for _, run := range sarif.Runs {
		run.DedupeArtifacts()
	}
	marshal, err := json.Marshal(sarif)
	if err != nil {
		return err
	}
	_, err = w.Write(marshal)
	return err
}

// PrettyWrite writes the JSON output with indentation
func (sarif *Report) PrettyWrite(w io.Writer) error {
	marshal, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(marshal)
	return err
}
