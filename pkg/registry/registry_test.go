package registry

import (
	"testing"
)

func TestSchemaPath(t *testing.T) {
	for i, testCase := range []struct {
		resourceKind, resourceAPIVersion, k8sVersion, expected string
		strict                                                 bool
	}{
		{
			"Deployment",
			"apps/v1",
			"1.16.0",
			"v1.16.0-standalone-strict/deployment-apps-v1.json",
			true,
		},
		{
			"Deployment",
			"apps/v1",
			"1.16.0",
			"v1.16.0-standalone/deployment-apps-v1.json",
			false,
		},
		{
			"Service",
			"v1",
			"1.18.0",
			"v1.18.0-standalone/service-v1.json",
			false,
		},
		{
			"Service",
			"v1",
			"master",
			"master-standalone/service-v1.json",
			false,
		},
	} {
		if got := schemaPath(testCase.resourceKind, testCase.resourceAPIVersion, testCase.k8sVersion, testCase.strict); got != testCase.expected {
			t.Errorf("%d - got %s, expected %s", i+1, got, testCase.expected)
		}
	}
}
