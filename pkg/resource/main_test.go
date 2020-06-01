package resource_test

import (
	"github.com/yannh/kubeconform/pkg/resource"
	"testing"
)

func TestSignatureFromBytes(t *testing.T) {
	testCases := []struct {
		name string
		have []byte
		want resource.Signature
		err error
	} {
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
			want: resource.Signature {
					Kind: "Deployment",
					Version: "apps/v1",
					Namespace: "default",
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		sig, err := resource.SignatureFromBytes(testCase.have)
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