package output

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/yannh/kubeconform/pkg/resource"

	"github.com/yannh/kubeconform/pkg/validator"
)

func TestJunitWrite(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		withSummary bool
		isStdin     bool
		verbose     bool
		results     []validator.Result
		expect      string
	}{
		{
			"an empty result",
			true,
			false,
			false,
			[]validator.Result{},
			"<testsuites name=\"kubeconform\" time=\"\" tests=\"0\" failures=\"0\" disabled=\"0\" errors=\"0\"></testsuites>\n",
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
			"<testsuites name=\"kubeconform\" time=\"\" tests=\"1\" failures=\"0\" disabled=\"0\" errors=\"0\">\n" +
				"  <testsuite name=\"deployment.yml\" id=\"1\" tests=\"1\" failures=\"0\" errors=\"0\" disabled=\"0\" skipped=\"0\">\n" +
				"    <properties></properties>\n" +
				"    <testcase name=\"my-app\" classname=\"Deployment@apps/v1\"></testcase>\n" +
				"  </testsuite>\n" +
				"</testsuites>\n",
		},
		{
			"a deployment, an empty resource, summary, no verbose",
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
				{
					Resource: resource.Resource{
						Path:  "deployment.yml",
						Bytes: []byte(`#A single comment`),
					},
					Status: validator.Empty,
					Err:    nil,
				},
			},
			"<testsuites name=\"kubeconform\" time=\"\" tests=\"1\" failures=\"0\" disabled=\"0\" errors=\"0\">\n" +
				"  <testsuite name=\"deployment.yml\" id=\"1\" tests=\"1\" failures=\"0\" errors=\"0\" disabled=\"0\" skipped=\"0\">\n" +
				"    <properties></properties>\n" +
				"    <testcase name=\"my-app\" classname=\"Deployment@apps/v1\"></testcase>\n" +
				"  </testsuite>\n" +
				"</testsuites>\n",
		},
	} {
		w := new(bytes.Buffer)
		o := junitOutput(w, testCase.withSummary, testCase.isStdin, testCase.verbose)

		for _, res := range testCase.results {
			o.Write(res)
		}
		o.Flush()

		// We remove the time, which will be different every time, before the comparison
		output := w.String()
		r := regexp.MustCompile(`time="[^"]*"`)
		output = r.ReplaceAllString(output, "time=\"\"")

		if output != testCase.expect {
			t.Errorf("%s - expected:, got:\n%s\n%s", testCase.name, testCase.expect, output)
		}
	}
}
