package validator

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/yannh/kubeconform/pkg/registry"

	"github.com/yannh/kubeconform/pkg/resource"
)

type mockRegistry struct {
	SchemaDownloader func() (string, []byte, error)
}

func newMockRegistry(f func() (string, []byte, error)) *mockRegistry {
	return &mockRegistry{
		SchemaDownloader: f,
	}
}

func (m mockRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (string, []byte, error) {
	return m.SchemaDownloader()
}

func TestValidate(t *testing.T) {
	for i, testCase := range []struct {
		name                         string
		rawResource, schemaRegistry1 []byte
		schemaRegistry2              []byte
		ignoreMissingSchema          bool
		strict                       bool
		injectDefaults               bool // New flag for default injection
		expectStatus                 Status
		expectErrors                 []ValidationError
		expectedInjectedDefaults     []string // Expected injected defaults
	}{
		{
			"valid resource with injected defaults",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
    "firstName": {
      "type": "string"
    },
    "lastName": {
      "type": "string",
      "default": "default_last_name"
    },
    "age": {
      "description": "Age in years",
      "type": "integer",
      "minimum": 0,
      "default": 30
    }
  },
  "required": ["firstName", "lastName"]
}`),
			nil,
			false,
			false,
			true, // Inject defaults
			Valid,
			[]ValidationError{},
			[]string{"lastName: default_last_name", "age: 30"}, // Expected injected defaults
		},
		{
			"invalid resource",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
    "firstName": {
      "type": "number"
    },
    "lastName": {
      "type": "string"
    },
    "age": {
      "description": "Age in years",
      "type": "integer",
      "minimum": 0
    }
  },
  "required": ["firstName", "lastName"]
}`),
			nil,
			false,
			false,
			false,
			Invalid,
			[]ValidationError{
				{
					Path: "/firstName",
					Msg:  "expected number, but got string",
				},
			},
			nil, // No defaults expected
		},
		{
			"missing required field",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
    "firstName": {
      "type": "string"
    },
    "lastName": {
      "type": "string"
    },
    "age": {
      "description": "Age in years",
      "type": "integer",
      "minimum": 0
    }
  },
  "required": ["firstName", "lastName"]
}`),
			nil,
			false,
			false,
			false,
			Invalid,
			[]ValidationError{
				{
					Path: "",
					Msg:  "missing properties: 'lastName'",
				},
			},
			nil, // No defaults expected
		},
	} {
		val := v{
			opts: Opts{
				SkipKinds:            map[string]struct{}{},
				RejectKinds:          map[string]struct{}{},
				IgnoreMissingSchemas: testCase.ignoreMissingSchema,
				Strict:               testCase.strict,
				InjectDefaults:       testCase.injectDefaults, // Inject defaults flag
			},
			schemaCache:    nil,
			schemaDownload: downloadSchema,
			regs: []registry.Registry{
				newMockRegistry(func() (string, []byte, error) {
					return "", testCase.schemaRegistry1, nil
				}),
				newMockRegistry(func() (string, []byte, error) {
					return "", testCase.schemaRegistry2, nil
				}),
			},
		}
		got := val.ValidateResource(resource.Resource{Bytes: testCase.rawResource})
		if got.Status != testCase.expectStatus {
			t.Errorf("Test '%s'- %d - expected %d, got %d", testCase.name, i, testCase.expectStatus, got.Status)
		}

		if len(got.ValidationErrors) != len(testCase.expectErrors) {
			t.Errorf("Test '%s': expected ValidationErrors: %+v, got: %v", testCase.name, testCase.expectErrors, got.ValidationErrors)
		}
		for i := range testCase.expectErrors {
			if testCase.expectErrors[i] != got.ValidationErrors[i] {
				t.Errorf("Test '%s': expected ValidationErrors: %+v, got: %v", testCase.name, testCase.expectErrors, got.ValidationErrors)
			}
		}

		// Check for injected defaults
		if testCase.injectDefaults && !reflect.DeepEqual(got.InjectedDefaults, testCase.expectedInjectedDefaults) {
			t.Errorf("Test '%s': expected injected defaults: %+v, got: %+v", testCase.name, testCase.expectedInjectedDefaults, got.InjectedDefaults)
		}
	}
}

func TestValidationErrors(t *testing.T) {
	rawResource := []byte(`
kind: name
apiVersion: v1
firstName: foo
age: not a number
`)

	schema := []byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
    "firstName": {
      "type": "string"
    },
    "lastName": {
      "type": "string"
    },
    "age": {
      "description": "Age in years",
      "type": "integer",
      "minimum": 0
    }
  },
  "required": ["firstName", "lastName"]
}`)

	expectedErrors := []ValidationError{
		{Path: "", Msg: "missing properties: 'lastName'"},
		{Path: "/age", Msg: "expected integer, but got string"},
	}

	val := v{
		opts: Opts{
			SkipKinds:   map[string]struct{}{},
			RejectKinds: map[string]struct{}{},
		},
		schemaCache:    nil,
		schemaDownload: downloadSchema,
		regs: []registry.Registry{
			newMockRegistry(func() (string, []byte, error) {
				return "", schema, nil
			}),
		},
	}

	got := val.ValidateResource(resource.Resource{Bytes: rawResource})
	if !reflect.DeepEqual(expectedErrors, got.ValidationErrors) {
		t.Errorf("Expected %+v, got %+v", expectedErrors, got.ValidationErrors)
	}
}

func TestValidateFile(t *testing.T) {
	inputData := []byte(`
kind: name
apiVersion: v1
firstName: bar
lastName: qux
---
kind: name
apiVersion: v1
firstName: foo
`)

	schema := []byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
    "firstName": {
      "type": "string"
    },
    "lastName": {
      "type": "string"
    }
  },
  "required": ["firstName", "lastName"]
}`)

	val := v{
		opts: Opts{
			SkipKinds:   map[string]struct{}{},
			RejectKinds: map[string]struct{}{},
		},
		schemaCache:    nil,
		schemaDownload: downloadSchema,
		regs: []registry.Registry{
			newMockRegistry(func() (string, []byte, error) {
				return "", schema, nil
			}),
		},
	}

	gotStatuses := []Status{}
	gotValidationErrors := []ValidationError{}
	for _, got := range val.Validate("test-file", io.NopCloser(bytes.NewReader(inputData))) {
		gotStatuses = append(gotStatuses, got.Status)
		gotValidationErrors = append(gotValidationErrors, got.ValidationErrors...)
	}

	expectedStatuses := []Status{Valid, Invalid}
	expectedValidationErrors := []ValidationError{
		{Path: "", Msg: "missing properties: 'lastName'"},
	}
	if !reflect.DeepEqual(expectedStatuses, gotStatuses) {
		t.Errorf("Expected %+v, got %+v", expectedStatuses, gotStatuses)
	}
	if !reflect.DeepEqual(expectedValidationErrors, gotValidationErrors) {
		t.Errorf("Expected %+v, got %+v", expectedValidationErrors, gotValidationErrors)
	}
}
