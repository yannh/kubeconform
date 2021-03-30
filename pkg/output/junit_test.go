package output

import (
	"bytes"
	"github.com/yannh/kubeconform/pkg/resource"
	"regexp"
	"testing"

	"github.com/yannh/kubeconform/pkg/validator"

	"github.com/beevik/etree"
)

func isNumeric(s string) bool {
	matched, _ := regexp.MatchString("^\\d+(\\.\\d+)?$", s)
	return matched
}

func TestJUnitWrite(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		withSummary bool
		isStdin     bool
		verbose     bool
		results     []validator.Result
		evaluate    func(d *etree.Document)
	}{
		{
			"empty document",
			false,
			false,
			false,
			[]validator.Result{},
			func(d *etree.Document) {
				root := d.FindElement("/testsuites")
				if root == nil {
					t.Errorf("Can't find root testsuite element")
					return
				}
				for _, attr := range root.Attr {
					switch attr.Key {
					case "time":
					case "tests":
					case "failures":
					case "disabled":
					case "errors":
						if !isNumeric(attr.Value) {
							t.Errorf("Expected a number for /testsuites/@%s", attr.Key)
						}
						continue
					case "name":
						if attr.Value != "kubeconform" {
							t.Errorf("Expected 'kubeconform' for /testsuites/@name")
						}
						continue
					default:
						t.Errorf("Unknown attribute /testsuites/@%s", attr.Key)
						continue
					}
				}
				suites := root.SelectElements("testsuite")
				if len(suites) != 0 {
					t.Errorf("No testsuite elements should be generated when there are no resources")
				}
			},
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
  namespace: "my-namespace"
`),
					},
					Status: validator.Valid,
					Err:    nil,
				},
			},
			func(d *etree.Document) {
				suites := d.FindElements("//testsuites/testsuite")
				if len(suites) != 1 {
					t.Errorf("Expected exactly 1 testsuite element, got %d", len(suites))
					return
				}
				suite := suites[0]
				for _, attr := range suite.Attr {
					switch attr.Key {
					case "name":
						if attr.Value != "deployment.yml" {
							t.Errorf("Test suite name should be the resource path")
						}
						continue
					case "tests":
						if attr.Value != "1" {
							t.Errorf("testsuite/@tests should be 1")
						}
						continue
					case "failures":
						if attr.Value != "0" {
							t.Errorf("testsuite/@failures should be 0")
						}
						continue
					case "errors":
						if attr.Value != "0" {
							t.Errorf("testsuite/@errors should be 0")
						}
						continue
					case "disabled":
						if attr.Value != "0" {
							t.Errorf("testsuite/@disabled should be 0")
						}
						continue
					case "skipped":
						if attr.Value != "0" {
							t.Errorf("testsuite/@skipped should be 0")
						}
						continue
					default:
						t.Errorf("Unknown testsuite attribute %s", attr.Key)
						continue
					}
				}
				testcases := suite.SelectElements("testcase")
				if len(testcases) != 1 {
					t.Errorf("Expected exactly 1 testcase, got %d", len(testcases))
					return
				}
				testcase := testcases[0]
				if testcase.SelectAttrValue("name", "") != "my-namespace/my-app" {
					t.Errorf("Test case name should be namespace / name")
				}
				if testcase.SelectAttrValue("classname", "") != "Deployment@apps/v1" {
					t.Errorf("Test case class name should be resource kind @ api version")
				}
				if testcase.SelectElement("skipped") != nil {
					t.Errorf("skipped element should not be generated if the kind was not skipped")
				}
				if testcase.SelectElement("error") != nil {
					t.Errorf("error element should not be generated if there was no error")
				}
				if len(testcase.SelectElements("failure")) != 0 {
					t.Errorf("failure elements should not be generated if there were no failures")
				}
			},
		},
	} {
		w := new(bytes.Buffer)
		o := junitOutput(w, testCase.withSummary, testCase.isStdin, testCase.verbose)

		for _, res := range testCase.results {
			o.Write(res)
		}
		o.Flush()

		doc := etree.NewDocument()
		doc.ReadFromString(w.String())

		testCase.evaluate(doc)
	}
}
