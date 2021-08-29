package output

import (
	"bytes"
	"testing"

	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

func TestTextWrite(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		withSummary bool
		isStdin     bool
		verbose     bool
		results     []validator.Result
		expect      string
	}{
		{
			"a single deployment, no summary, no verbose",
			false,
			false,
			false,
			[]validator.Result{},
			"",
		},
		{
			"a single deployment, summary, no verbose",
			true,
			false,
			false,
			[]validator.Result{
				{
					Resource: resource.Resource{
						Path: "deployment.yml",
						Bytes: []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: "my-app"
`),
					},
					Status: validator.Valid,
					Err:    nil,
				},
			},
			"Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0\n",
		},
		{
			"a single deployment, verbose, with summary",
			true,
			false,
			true,
			[]validator.Result{
				{
					Resource: resource.Resource{
						Path: "deployment.yml",
						Bytes: []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: "my-app"
`),
					},
					Status: validator.Valid,
					Err:    nil,
				},
			},
			`deployment.yml - Deployment my-app is valid
Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0
`,
		},
	} {
		w := new(bytes.Buffer)
		o := textOutput(w, testCase.withSummary, testCase.isStdin, testCase.verbose)

		for _, res := range testCase.results {
			o.Write(res)
		}
		o.Flush()

		if w.String() != testCase.expect {
			t.Errorf("%s - expected: %s, got: %s", testCase.name, testCase.expect, w)
		}
	}
}
