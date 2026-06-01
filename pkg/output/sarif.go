package output

import (
	"fmt"
	"io"
	"sync"

	"github.com/owenrumney/go-sarif/v3/pkg/report"
	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

const (
	toolName    = "kubeconform"
	toolInfoURI = "https://github.com/yannh/kubeconform"
)

const (
	ruleIDValid   = "KUBE-VALID"
	ruleIDInvalid = "KUBE-INVALID"
	ruleIDError   = "KUBE-ERROR"
	ruleIDSkipped = "KUBE-SKIPPED"
)

const (
	levelNote  = "note"
	levelError = "error"
)

var sarifReportingDescriptors = []*sarif.ReportingDescriptor{
	newSarifReportingDescriptor(ruleIDValid, "ValidResource", "Resource is valid.", levelNote),
	newSarifReportingDescriptor(ruleIDInvalid, "InvalidResource", "Resource is invalid against schema.", levelError),
	newSarifReportingDescriptor(ruleIDError, "ProcessingError", "Error processing resource.", levelError),
	newSarifReportingDescriptor(ruleIDSkipped, "SkippedResource", "Resource validation was skipped.", levelNote),
}

// newSarifReportingDescriptor creates a new SARIF reporting descriptor.
func newSarifReportingDescriptor(id, name, shortDesc, level string) *sarif.ReportingDescriptor {
	shortDescMsg := sarif.NewMultiformatMessageString().WithText(shortDesc)

	return sarif.NewRule(id).
		WithName(name).
		WithShortDescription(shortDescMsg).
		WithDefaultConfiguration(sarif.NewReportingConfiguration().WithLevel(level))
}

// sarifOutputter handles the generation of SARIF format output.
// It implements the Output interface and is concurrency-safe.
type sarifOutputter struct {
	mu      sync.Mutex
	writer  io.Writer
	verbose bool
	results []*sarif.Result
}

// sarifOutput creates a new Outputter that formats results as SARIF.
func sarifOutput(writer io.Writer, withSummary, isStdin, verbose bool) Output {
	return &sarifOutputter{
		writer:  writer,
		verbose: verbose,
		results: make([]*sarif.Result, 0),
	}
}

// newSarifRun creates and initializes a new SARIF report and run
// with the standard tool and rule information.
func newSarifRun() (*sarif.Report, *sarif.Run, error) {
	rep := report.NewV210Report()
	if rep == nil {
		return nil, nil, fmt.Errorf("failed to initialize SARIF report")
	}

	run := sarif.NewRunWithInformationURI(toolName, toolInfoURI)
	if run == nil {
		return nil, nil, fmt.Errorf("failed to initialize SARIF run")
	}

	if run.Tool == nil || run.Tool.Driver == nil {
		return nil, nil, fmt.Errorf("SARIF run is missing required tool driver information")
	}

	run.Tool.Driver.WithRules(sarifReportingDescriptors)
	rep.AddRun(run)

	return rep, run, nil
}

// Write processes a single validation result.
// It is concurrency-safe.
func (so *sarifOutputter) Write(validationResult validator.Result) error {
	so.mu.Lock()
	defer so.mu.Unlock()

	if validationResult.Status == validator.Empty {
		return nil
	}

	if validationResult.Status == validator.Valid && !so.verbose {
		return nil
	}

	signature, _ := validationResult.Resource.Signature()

	if validationResult.Status == validator.Invalid && len(validationResult.ValidationErrors) > 0 {
		for _, valErr := range validationResult.ValidationErrors {
			sarifResult := so.newSarifResult(validationResult, signature, &valErr)
			so.results = append(so.results, sarifResult)
		}
	} else {
		sarifResult := so.newSarifResult(validationResult, signature, nil)
		so.results = append(so.results, sarifResult)
	}
	return nil
}

// newSarifResult creates a SARIF result from a validation result.
// If valErr is provided, it populates the result with specific validation
// failure details, including the logical path.
func (so *sarifOutputter) newSarifResult(res validator.Result, sig *resource.Signature, valErr *validator.ValidationError) *sarif.Result {
	result := sarif.NewResult().
		AddLocation(
			sarif.NewLocationWithPhysicalLocation(
				sarif.NewPhysicalLocation().
					WithArtifactLocation(
						sarif.NewSimpleArtifactLocation(res.Resource.Path),
					),
			),
		)

	if valErr != nil {
		result.Locations[0].AddLogicalLocation(
			sarif.NewLogicalLocation().WithName(valErr.Path),
		)
		result.
			WithRuleID(ruleIDInvalid).
			WithLevel(levelError)

		var message string
		if sig.Kind != "" && sig.Name != "" {
			message = fmt.Sprintf("%s %s is invalid: %s: %s", sig.Kind, sig.Name, valErr.Path, valErr.Msg)
		} else {
			message = fmt.Sprintf("%s is invalid: %s: %s", res.Resource.Path, valErr.Path, valErr.Msg)
		}
		result.WithMessage(sarif.NewTextMessage(message))
	} else {
		switch res.Status {
		case validator.Valid:
			result.
				WithRuleID(ruleIDValid).
				WithLevel(levelNote).
				WithMessage(sarif.NewTextMessage(
					fmt.Sprintf("%s %s is valid", sig.Kind, sig.Name),
				))
		case validator.Error:
			result.
				WithRuleID(ruleIDError).
				WithLevel(levelError)
			var message string
			if sig.Kind != "" && sig.Name != "" {
				message = fmt.Sprintf("%s %s failed validation: %s", sig.Kind, sig.Name, res.Err.Error())
			} else {
				message = fmt.Sprintf("%s failed validation: %s", res.Resource.Path, res.Err.Error())
			}
			result.WithMessage(sarif.NewTextMessage(message))
		case validator.Skipped:
			result.
				WithRuleID(ruleIDSkipped).
				WithLevel(levelNote).
				WithMessage(sarif.NewTextMessage(
					fmt.Sprintf("%s %s skipped", sig.Kind, sig.Name),
				))
		default:
			result.
				WithRuleID(ruleIDError).
				WithLevel(levelError).
				WithMessage(sarif.NewTextMessage(
					fmt.Sprintf("Unknown validation status for %s", res.Resource.Path),
				))
		}
	}
	return result
}

// Flush generates the complete SARIF report and writes it to the output writer.
// It is concurrency-safe.
func (so *sarifOutputter) Flush() error {
	so.mu.Lock()
	defer so.mu.Unlock()

	rep, run, err := newSarifRun()
	if err != nil {
		return err
	}

	for _, result := range so.results {
		run.AddResult(result)
	}

	if err := rep.PrettyWrite(so.writer); err != nil {
		return fmt.Errorf("failed to write SARIF report: %w", err)
	}

	return nil
}
