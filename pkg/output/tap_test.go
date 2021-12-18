package output

import (
	"bytes"
	"testing"

	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

func TestTapWrite(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		withSummary bool
		isStdin     bool
		verbose     bool
		results     []validator.Result
		expect      string
	}{
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
			"TAP version 13\nok 1 - deployment.yml (apps/v1/Deployment//my-app)\n1..1\n",
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
			"TAP version 13\nok 1 - deployment.yml (apps/v1/Deployment//my-app)\n1..1\n",
		},
	} {
		w := new(bytes.Buffer)
		o := tapOutput(w, testCase.withSummary, testCase.isStdin, testCase.verbose)

		for _, res := range testCase.results {
			o.Write(res)
		}
		o.Flush()

		if w.String() != testCase.expect {
			t.Errorf("%s - expected:, got:\n%s\n%s", testCase.name, testCase.expect, w)
		}
	}
}
