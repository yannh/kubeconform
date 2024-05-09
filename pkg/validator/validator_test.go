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
		expectStatus                 Status
		expectErrors                 []ValidationError
	}{
		{
			"valid resource",
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
			Valid,
			[]ValidationError{},
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
			Invalid,
			[]ValidationError{
				{
					Path: "/firstName",
					Msg:  "expected number, but got string",
				},
			},
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
			Invalid,
			[]ValidationError{
				{
					Path: "",
					Msg:  "missing properties: 'lastName'",
				},
			},
		},
		{
			"key \"firstName\" already set in map",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
firstName: bar
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
    }
  },
  "required": ["firstName"]
}`),
			nil,
			false,
			true,
			Error,
			[]ValidationError{},
		},
		{
			"key firstname already set in map in non-strict mode",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
firstName: bar
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
    }
  },
  "required": ["firstName"]
}`),
			nil,
			false,
			false,
			Valid,
			[]ValidationError{},
		},
		{
			"resource has invalid yaml",
			[]byte(`
kind: name
apiVersion: v1
firstName foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
    "apiVersion": {
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
			Error,
			[]ValidationError{},
		},
		{
			"missing schema in 1st registry",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
lastName: bar
`),
			nil,
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
    "apiVersion": {
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
			false,
			false,
			Valid,
			[]ValidationError{},
		},
		{
			"non-json response in 1st registry",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
lastName: bar
`),
			[]byte(`<html>error page</html>`),
			[]byte(`{
 "title": "Example Schema",
 "type": "object",
 "properties": {
   "kind": {
     "type": "string"
   },
   "apiVersion": {
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
			false,
			false,
			Valid,
			[]ValidationError{},
		},
		{
			"missing schema in both registries, ignore missing",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
lastName: bar
`),
			nil,
			nil,
			true,
			false,
			Skipped,
			[]ValidationError{},
		},
		{
			"missing schema in both registries, do not ignore missing",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
lastName: bar
`),
			nil,
			nil,
			false,
			false,
			Error,
			[]ValidationError{},
		},
		{
			"non-json response in both registries, ignore missing",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
lastName: bar
`),
			[]byte(`<html>error page</html>`),
			[]byte(`<html>error page</html>`),
			true,
			false,
			Skipped,
			[]ValidationError{},
		},
		{
			"non-json response in both registries, do not ignore missing",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
lastName: bar
`),
			[]byte(`<html>error page</html>`),
			[]byte(`<html>error page</html>`),
			false,
			false,
			Error,
			[]ValidationError{},
		},
	} {
		val := v{
			opts: Opts{
				SkipKinds:            map[string]struct{}{},
				RejectKinds:          map[string]struct{}{},
				IgnoreMissingSchemas: testCase.ignoreMissingSchema,
				Strict:               testCase.strict,
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
			if got.Err != nil {
				t.Errorf("Test '%s' - expected %d, got %d: %s", testCase.name, testCase.expectStatus, got.Status, got.Err.Error())
			} else {
				t.Errorf("Test '%s'- %d - expected %d, got %d", testCase.name, i, testCase.expectStatus, got.Status)
			}
		}

		if len(got.ValidationErrors) != len(testCase.expectErrors) {
			t.Errorf("Test '%s': expected ValidationErrors: %+v, got: % v", testCase.name, testCase.expectErrors, got.ValidationErrors)
		}
		for i, _ := range testCase.expectErrors {
			if testCase.expectErrors[i] != got.ValidationErrors[i] {
				t.Errorf("Test '%s': expected ValidationErrors: %+v, got: % v", testCase.name, testCase.expectErrors, got.ValidationErrors)
			}
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
