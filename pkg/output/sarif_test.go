package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
	"gopkg.in/yaml.v2"
)

type testMetadata struct {
	Name string `yaml:"name"`
}
type testResource struct {
	APIVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   testMetadata `yaml:"metadata"`
}

func marshalTestResource(t *testing.T, apiVersion, kind, name string) []byte {
	res := testResource{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata:   testMetadata{Name: name},
	}
	bytes, err := yaml.Marshal(res)
	require.NoError(t, err)
	return bytes
}

// newExpectedReport creates a complete, initialized SARIF report
// for use in test assertions.
func newExpectedReport(t *testing.T, results []*sarif.Result) *sarif.Report {
	report, run, err := newSarifRun()
	require.NoError(t, err)

	if results != nil {
		run.Results = results
	}

	report.InlineExternalProperties = nil
	report.Properties.Tags = nil
	run.Artifacts = nil

	for _, rule := range run.Tool.Driver.Rules {
		rule.DeprecatedGuids = nil
		rule.DeprecatedIds = nil
		rule.DeprecatedNames = nil
	}

	return report
}

func newExpectedResult(ruleID, level, message, path string, logicalPath string) *sarif.Result {
	result := sarif.NewResult().
		WithRuleID(ruleID).
		WithLevel(level).
		WithMessage(sarif.NewTextMessage(message))

	location := sarif.NewLocationWithPhysicalLocation(
		sarif.NewPhysicalLocation().
			WithArtifactLocation(
				sarif.NewSimpleArtifactLocation(path),
			),
	)

	if logicalPath != "" {
		location.AddLogicalLocation(sarif.NewLogicalLocation().WithName(logicalPath))
	}
	result.AddLocation(location)

	result.Suppressions = nil
	result.WorkItemUris = nil

	return result
}

func TestSarifWrite(t *testing.T) {
	testCases := []struct {
		name           string
		verbose        bool
		results        []validator.Result
		expectedReport *sarif.Report
	}{
		{
			name:    "single invalid deployment",
			verbose: false,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "deployment.yml",
						Bytes: marshalTestResource(t, "apps/v1", "Deployment", "my-app"),
					},
					Status: validator.Invalid,
					Err:    errors.New("spec.replicas: Invalid type. Expected: [integer,null], given: string"),
					ValidationErrors: []validator.ValidationError{
						{Path: "spec.replicas", Msg: "Invalid type. Expected: [integer,null], given: string"},
					},
				},
			},
			expectedReport: newExpectedReport(t, []*sarif.Result{
				newExpectedResult(
					ruleIDInvalid,
					levelError,
					"Deployment my-app is invalid: spec.replicas: Invalid type. Expected: [integer,null], given: string",
					"deployment.yml",
					"spec.replicas",
				),
			}),
		},
		{
			name:    "single valid deployment verbose",
			verbose: true,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "deployment.yml",
						Bytes: marshalTestResource(t, "apps/v1", "Deployment", "my-app"),
					},
					Status: validator.Valid,
					Err:    nil,
				},
			},
			expectedReport: newExpectedReport(t, []*sarif.Result{
				newExpectedResult(
					ruleIDValid,
					levelNote,
					"Deployment my-app is valid",
					"deployment.yml",
					"",
				),
			}),
		},
		{
			name:    "single valid deployment non-verbose",
			verbose: false,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "deployment.yml",
						Bytes: marshalTestResource(t, "apps/v1", "Deployment", "my-app"),
					},
					Status: validator.Valid,
					Err:    nil,
				},
			},
			expectedReport: newExpectedReport(t, nil),
		},
		{
			name:    "skipped resource",
			verbose: true,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "service.yml",
						Bytes: marshalTestResource(t, "v1", "Service", "my-service"),
					},
					Status: validator.Skipped,
					Err:    nil,
				},
			},
			expectedReport: newExpectedReport(t, []*sarif.Result{
				newExpectedResult(
					ruleIDSkipped,
					levelNote,
					"Service my-service skipped",
					"service.yml",
					"",
				),
			}),
		},
		{
			name:    "error processing resource",
			verbose: false,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "configmap.yml",
						Bytes: marshalTestResource(t, "v1", "ConfigMap", "my-config"),
					},
					Status: validator.Error,
					Err:    errors.New("failed to download schema"),
				},
			},
			expectedReport: newExpectedReport(t, []*sarif.Result{
				newExpectedResult(
					ruleIDError,
					levelError,
					"ConfigMap my-config failed validation: failed to download schema",
					"configmap.yml",
					"",
				),
			}),
		},
		{
			name:    "empty resource filtered out",
			verbose: true,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "empty.yml",
						Bytes: []byte(`---`),
					},
					Status: validator.Empty,
					Err:    nil,
				},
			},
			expectedReport: newExpectedReport(t, nil),
		},
		{
			name:    "multiple invalid results from one file",
			verbose: false,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "deployment1.yml",
						Bytes: marshalTestResource(t, "apps/v1", "Deployment", "app1"),
					},
					Status: validator.Invalid,
					ValidationErrors: []validator.ValidationError{
						{Path: "spec.template", Msg: "is missing"},
						{Path: "spec.selector", Msg: "is missing"},
					},
				},
				{
					Resource: resource.Resource{
						Path:  "deployment2.yml",
						Bytes: marshalTestResource(t, "apps/v1", "Deployment", "app2"),
					},
					Status: validator.Invalid,
					ValidationErrors: []validator.ValidationError{
						{Path: "spec.replicas", Msg: "must be positive"},
					},
				},
			},
			expectedReport: newExpectedReport(t, []*sarif.Result{
				newExpectedResult(
					ruleIDInvalid,
					levelError,
					"Deployment app1 is invalid: spec.template: is missing",
					"deployment1.yml",
					"spec.template",
				),
				newExpectedResult(
					ruleIDInvalid,
					levelError,
					"Deployment app1 is invalid: spec.selector: is missing",
					"deployment1.yml",
					"spec.selector",
				),
				newExpectedResult(
					ruleIDInvalid,
					levelError,
					"Deployment app2 is invalid: spec.replicas: must be positive",
					"deployment2.yml",
					"spec.replicas",
				),
			}),
		},
		{
			name:    "invalid resource with no signature",
			verbose: false,
			results: []validator.Result{
				{
					Resource: resource.Resource{
						Path:  "not-yaml.yml",
						Bytes: []byte(`not: valid: yaml`),
					},
					Status: validator.Invalid,
					ValidationErrors: []validator.ValidationError{
						{Path: "metadata", Msg: "is missing"},
					},
				},
			},
			expectedReport: newExpectedReport(t, []*sarif.Result{
				newExpectedResult(
					ruleIDInvalid,
					levelError,
					"not-yaml.yml is invalid: metadata: is missing",
					"not-yaml.yml",
					"metadata",
				),
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := new(bytes.Buffer)
			o := sarifOutput(w, false, false, tc.verbose)

			for _, res := range tc.results {
				err := o.Write(res)
				require.NoError(t, err, "Write() should not return an error")
			}

			err := o.Flush()
			require.NoError(t, err, "Flush() should not return an error")

			var actualReport sarif.Report
			err = json.Unmarshal(w.Bytes(), &actualReport)
			require.NoError(t, err, "Output should be valid JSON: %s", w.String())

			require.Len(t, actualReport.Runs, 1)
			assert.Equal(t, tc.expectedReport.Version, actualReport.Version)
			assert.Equal(t, tc.expectedReport.Schema, actualReport.Schema)
			assert.Equal(t, tc.expectedReport.Properties, actualReport.Properties)
			assert.Equal(t, tc.expectedReport.Runs[0].Tool, actualReport.Runs[0].Tool)

			assert.ElementsMatch(t, tc.expectedReport.Runs[0].Results, actualReport.Runs[0].Results)
		})
	}
}
