package registry

import (
	"testing"
)

func TestSchemaPath(t *testing.T) {
	for i, testCase := range []struct {
		tpl, resourceKind, resourceAPIVersion, k8sVersion, expected string
		strict                                                      bool
		delims                                                      string
		errExpected                                                 error
	}{
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Deployment",
			"apps/v1",
			"1.16.0",
			"https://kubernetesjsonschema.dev/v1.16.0-standalone-strict/deployment-apps-v1.json",
			true,
			"",
			nil,
		},
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Deployment",
			"apps/v1",
			"1.16.0",
			"https://kubernetesjsonschema.dev/v1.16.0-standalone/deployment-apps-v1.json",
			false,
			"",
			nil,
		},
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Service",
			"v1",
			"1.18.0",
			"https://kubernetesjsonschema.dev/v1.18.0-standalone/service-v1.json",
			false,
			"",
			nil,
		},
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Service",
			"v1",
			"master",
			"https://kubernetesjsonschema.dev/master-standalone/service-v1.json",
			false,
			"",
			nil,
		},
		{
			"https://kubernetesjsonschema.dev/[[ .NormalizedKubernetesVersion ]]-standalone[[ .StrictSuffix ]]/[[ .ResourceKind ]][[ .KindSuffix ]].json",
			"Service",
			"v1",
			"master",
			"https://kubernetesjsonschema.dev/master-standalone/service-v1.json",
			false,
			"[[,]]",
			nil,
		},
	} {
		got, err := schemaPath(testCase.tpl, testCase.resourceKind, testCase.resourceAPIVersion, testCase.k8sVersion, testCase.strict, testCase.delims)
		if err != testCase.errExpected {
			t.Errorf("%d - got error %s, expected %s", i+1, err, testCase.errExpected)
		}
		if got != testCase.expected {
			t.Errorf("%d - got %s, expected %s", i+1, got, testCase.expected)
		}
	}
}
