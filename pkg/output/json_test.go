package output

import (
	"bytes"
	"testing"

	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

func TestJSONWrite(t *testing.T) {
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
			"{\n  \"resources\": []\n}\n",
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
			`{
  "resources": [],
  "summary": {
    "valid": 1,
    "invalid": 0,
    "errors": 0,
    "skipped": 0
  }
}
`,
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
			`{
  "resources": [
    {
      "filename": "deployment.yml",
      "kind": "Deployment",
      "name": "my-app",
      "version": "apps/v1",
      "status": "statusValid",
      "msg": ""
    }
  ],
  "summary": {
    "valid": 1,
    "invalid": 0,
    "errors": 0,
    "skipped": 0
  }
}
`,
		},
		{
			"a single invalid deployment, verbose, with summary",
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
					Status: validator.Invalid,
					Err: &validator.ValidationError{
						Path: "foo",
						Msg:  "bar",
					},
					ValidationErrors: []validator.ValidationError{
						{
							Path: "foo",
							Msg:  "bar",
						},
					},
				},
			},
			`{
  "resources": [
    {
      "filename": "deployment.yml",
      "kind": "Deployment",
      "name": "my-app",
      "version": "apps/v1",
      "status": "statusInvalid",
      "msg": "bar",
      "validationErrors": [
        {
          "path": "foo",
          "msg": "bar"
        }
      ]
    }
  ],
  "summary": {
    "valid": 0,
    "invalid": 1,
    "errors": 0,
    "skipped": 0
  }
}
`,
		},
	} {
		w := new(bytes.Buffer)
		o := jsonOutput(w, testCase.withSummary, testCase.isStdin, testCase.verbose)

		for _, res := range testCase.results {
			o.Write(res)
		}
		o.Flush()

		if w.String() != testCase.expect {
			t.Errorf("%s - expected: %s, got: %s", testCase.name, testCase.expect, w)
		}
	}
}
