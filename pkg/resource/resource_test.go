package resource_test

import (
	"log"
	"reflect"
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/yannh/kubeconform/pkg/resource"
)

func TestSignatureFromBytes(t *testing.T) {
	testCases := []struct {
		name string
		have []byte
		want resource.Signature
		err  error
	}{
		{
			name: "valid deployment",
			have: []byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myService
  namespace: default
  labels:
    app: myService
spec:
`),
			want: resource.Signature{
				Kind:      "Deployment",
				Version:   "apps/v1",
				Namespace: "default",
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		res := resource.Resource{Bytes: testCase.have}
		sig, err := res.Signature()
		if err != nil && err.Error() != testCase.err.Error() {
			t.Errorf("test \"%s\" - received error: %s", testCase.name, err)
		}
		if sig.Version != testCase.want.Version ||
			sig.Kind != testCase.want.Kind ||
			sig.Namespace != testCase.want.Namespace {
			t.Errorf("test \"%s\": received %+v, expected %+v", testCase.name, sig, testCase.want)
		}
	}
}

func TestSignatureFromMap(t *testing.T) {
	testCases := []struct {
		b string
		s resource.Signature
	}{
		{
			"apiVersion: v1\nkind: ReplicationController\nmetadata:\n  name: \"bob\"\nspec:\n  replicas: 2\n",
			resource.Signature{
				Kind:      "ReplicationController",
				Version:   "v1",
				Namespace: "",
				Name:      "bob",
			},
		},
	}

	for i, testCase := range testCases {
		res := resource.Resource{
			Path:  "foo",
			Bytes: []byte(testCase.b),
		}

		var r map[string]interface{}
		if err := yaml.Unmarshal(res.Bytes, &r); err != nil {
			log.Fatal(err)
		}

		res.SignatureFromMap(r)
		sig, _ := res.Signature()
		if !reflect.DeepEqual(*sig, testCase.s) {
			t.Errorf("test %d - for resource %s, expected %+v, got %+v", i+1, testCase.b, testCase.s, sig)
		}
	}
}

func TestResources(t *testing.T) {
	testCases := []struct {
		b        string
		expected int
	}{
		{
			`
apiVersion: v1
kind: List
`,
			0,
		},
		{
			`
apiVersion: v1
kind: List
Items: []
`,
			0,
		},
		{
			`
apiVersion: v1
kind: List
Items:
- apiVersion: v1
  kind: ReplicationController
  metadata:
    name: "bob"
  spec:
    replicas: 2
`,
			1,
		},
		{
			`
apiVersion: v1
kind: List
Items:
- apiVersion: v1
  kind: ReplicationController
  metadata:
    name: "bob"
  spec:
    replicas: 2
- apiVersion: v1
  kind: ReplicationController
  metadata:
    name: "Jim"
  spec:
    replicas: 2
`,
			2,
		},
	}

	for i, testCase := range testCases {
		res := resource.Resource{
			Path:  "foo",
			Bytes: []byte(testCase.b),
		}

		subres := res.Resources()
		if len(subres) != testCase.expected {
			t.Errorf("test %d: expected to find %d resources, found %d", i, testCase.expected, len(subres))
		}
	}
}
