package resource_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/yannh/kubeconform/pkg/resource"
	"sigs.k8s.io/yaml"
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
	}{
		{
			"apiVersion: v1\nkind: ReplicationController\nmetadata:\n  name: \"bob\"\nspec:\n  replicas: 2\n",
		},
	}

	for _, testCase := range testCases {
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
		fmt.Printf("%+v", sig)
	}
}
