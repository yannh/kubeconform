package registry

import (
	"testing"
)

func TestSchemaPath(t *testing.T) {
	for i, testCase := range []struct {
		tpl, resourceKind, resourceAPIVersion, k8sVersion, expected string
		strict                                                      bool
		errExpected                                                 error
	}{
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Deployment",
			"apps/v1",
			"1.16.0",
			"https://kubernetesjsonschema.dev/v1.16.0-standalone-strict/deployment-apps-v1.json",
			true,
			nil,
		},
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Deployment",
			"apps/v1",
			"1.16.0",
			"https://kubernetesjsonschema.dev/v1.16.0-standalone/deployment-apps-v1.json",
			false,
			nil,
		},
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Service",
			"v1",
			"1.18.0",
			"https://kubernetesjsonschema.dev/v1.18.0-standalone/service-v1.json",
			false,
			nil,
		},
		{
			"https://kubernetesjsonschema.dev/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json",
			"Service",
			"v1",
			"master",
			"https://kubernetesjsonschema.dev/master-standalone/service-v1.json",
			false,
			nil,
		},
	} {
		got, err := schemaPath(testCase.tpl, testCase.resourceKind, testCase.resourceAPIVersion, testCase.k8sVersion, testCase.strict)
		if err != testCase.errExpected {
			t.Errorf("%d - got error %s, expected %s", i+1, err, testCase.errExpected)
		}
		if got != testCase.expected {
			t.Errorf("%d - got %s, expected %s", i+1, got, testCase.expected)
		}
	}
}

func TestNewSchemaLocationTemplate(t *testing.T) {
	defaultTpl := "https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json"
	defaultSuffix := "/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json"

	for i, testCase := range []struct {
		name     string
		location string
		want     string
	}{
		{
			name:     "default keyword",
			location: "default",
			want:     defaultTpl,
		},
		{
			name:     "bare directory gets default path appended",
			location: "https://example.com/schemas",
			want:     "https://example.com/schemas" + defaultSuffix,
		},
		{
			name:     "full templated path ending in json is unchanged",
			location: "https://example.com/{{ .ResourceKind }}.json",
			want:     "https://example.com/{{ .ResourceKind }}.json",
		},
		{
			// GitLab raw API puts /raw?... after the filename, so the template
			// does not end in "json" even though it already specifies the path.
			name:     "full templated path not ending in json is unchanged",
			location: "https://gitlab.example.com/api/v4/projects/123/repository/files/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}%2F{{ .ResourceKind }}{{ .KindSuffix }}.json/raw?ref=master",
			want:     "https://gitlab.example.com/api/v4/projects/123/repository/files/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}%2F{{ .ResourceKind }}{{ .KindSuffix }}.json/raw?ref=master",
		},
	} {
		reg, err := New(testCase.location, "", false, false, false)
		if err != nil {
			t.Fatalf("%d %s: New() error: %v", i+1, testCase.name, err)
		}
		got, ok := reg.(*SchemaRegistry)
		if !ok {
			t.Fatalf("%d %s: expected *SchemaRegistry, got %T", i+1, testCase.name, reg)
		}
		if got.schemaPathTemplate != testCase.want {
			t.Errorf("%d %s: schemaPathTemplate = %q, want %q", i+1, testCase.name, got.schemaPathTemplate, testCase.want)
		}
	}
}
