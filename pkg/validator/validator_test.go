package validator

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/yannh/kubeconform/pkg/loader"

	"github.com/yannh/kubeconform/pkg/registry"

	"github.com/yannh/kubeconform/pkg/resource"
)

type mockRegistry struct {
	SchemaDownloader func() (string, any, error)
}

func newMockRegistry(f func() (string, any, error)) *mockRegistry {
	return &mockRegistry{
		SchemaDownloader: f,
	}
}

func (m mockRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (string, any, error) {
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
					Msg:  "got string, want number",
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
					Msg:  "missing property 'lastName'",
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
		{
			"valid resource duration - go format",
			[]byte(`
kind: name
apiVersion: v1
interval: 5s
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
	"interval": {
      "type": "string",
	  "format": "duration"
    }
  },
  "required": ["interval"]
}`),
			nil,
			false,
			false,
			Valid,
			[]ValidationError{},
		},
		{
			"valid resource duration - scala duration format",
			[]byte(`
kind: name
apiVersion: v1
interval: 2w
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
	"interval": {
      "type": "string",
	  "format": "duration"
    }
  },
  "required": ["interval"]
}`),
			nil,
			false,
			false,
			Valid,
			[]ValidationError{},
		},
		{
			"valid resource duration - iso8601 format",
			[]byte(`
kind: name
apiVersion: v1
interval: PT1H
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
	"interval": {
      "type": "string",
	  "format": "duration"
    }
  },
  "required": ["interval"]
}`),
			nil,
			false,
			false,
			Valid,
			[]ValidationError{},
		},
		{
			"invalid resource duration",
			[]byte(`
kind: name
apiVersion: v1
interval: test
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "kind": {
      "type": "string"
    },
	"interval": {
      "type": "string",
	  "format": "duration"
    }
  },
  "required": ["interval"]
}`),
			nil,
			false,
			false,
			Invalid,
			[]ValidationError{{Path: "/interval", Msg: "'test' is not valid duration: must start with P"}},
		},
		{
			// A schema file whose entire content is the literal "null" decodes to
			// a nil document. It used to reach jsonschema's Compile and panic with
			// a nil pointer dereference (issue #337); it must now be treated as a
			// missing schema instead.
			"schema document is null",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
`),
			[]byte(`null`),
			nil,
			true,
			false,
			Skipped,
			[]ValidationError{},
		},
		{
			// Same as above but with IgnoreMissingSchemas disabled: a null schema
			// document must surface as a graceful Error, not a panic (issue #337).
			"schema document is null, do not ignore missing",
			[]byte(`
kind: name
apiVersion: v1
firstName: foo
`),
			[]byte(`null`),
			nil,
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
			schemaDownload: downloadSchema,
			regs: []registry.Registry{
				newMockRegistry(func() (string, any, error) {
					if testCase.schemaRegistry1 == nil {
						return "", nil, loader.NewNotFoundError(nil)
					}
					s, err := jsonschema.UnmarshalJSON(bytes.NewReader(testCase.schemaRegistry1))
					if err != nil {
						return "", s, loader.NewNonJSONResponseError(err)
					}
					return "", s, err
				}),
				newMockRegistry(func() (string, any, error) {
					if testCase.schemaRegistry2 == nil {
						return "", nil, loader.NewNotFoundError(nil)
					}
					s, err := jsonschema.UnmarshalJSON(bytes.NewReader(testCase.schemaRegistry2))
					if err != nil {
						return "", s, loader.NewNonJSONResponseError(err)
					}
					return "", s, err
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
		{Path: "", Msg: "missing property 'lastName'"},
		{Path: "/age", Msg: "got string, want integer"},
	}

	val := v{
		opts: Opts{
			SkipKinds:   map[string]struct{}{},
			RejectKinds: map[string]struct{}{},
		},
		schemaDownload: downloadSchema,
		regs: []registry.Registry{
			newMockRegistry(func() (string, any, error) {
				s, err := jsonschema.UnmarshalJSON(bytes.NewReader(schema))
				if err != nil {
					return "", s, loader.NewNonJSONResponseError(err)
				}
				return "", s, err
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
		schemaDownload: downloadSchema,
		regs: []registry.Registry{
			newMockRegistry(func() (string, any, error) {
				s, err := jsonschema.UnmarshalJSON(bytes.NewReader(schema))
				return "", s, err
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
		{Path: "", Msg: "missing property 'lastName'"},
	}
	if !reflect.DeepEqual(expectedStatuses, gotStatuses) {
		t.Errorf("Expected %+v, got %+v", expectedStatuses, gotStatuses)
	}
	if !reflect.DeepEqual(expectedValidationErrors, gotValidationErrors) {
		t.Errorf("Expected %+v, got %+v", expectedValidationErrors, gotValidationErrors)
	}
}

// --- issue #296 / #305: a schema that is served but cannot be compiled ---
// Refs https://github.com/yannh/kubeconform/issues/296
// Refs https://github.com/yannh/kubeconform/issues/305

// unusableSchemaRegistry serves a document that parses as JSON but is not a
// valid JSON schema: draft-04 requires "type" to be a string or an array.
func unusableSchemaRegistry() registry.Registry {
	return newMockRegistry(func() (string, any, error) {
		s, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(`{"type": 123}`)))
		return "unusable.json", s, err
	})
}

func notFoundRegistry() registry.Registry {
	return newMockRegistry(func() (string, any, error) {
		return "", nil, loader.NewNotFoundError(nil)
	})
}

// RED: downloadSchema must not report "no schema anywhere" when a registry did
// serve a schema and compiling it failed.
func TestDownloadSchemaReportsUnusableSchema(t *testing.T) {
	schema, err := downloadSchema(
		[]registry.Registry{unusableSchemaRegistry()},
		jsonschema.SchemeURLLoader{}, "name", "v1", "1.30.0")

	if schema != nil {
		t.Fatalf("expected no schema, got %v", schema)
	}
	if err == nil {
		t.Fatalf("expected the compile error to be returned, got nil - the reason the schema could not be used is lost")
	}
	if !strings.Contains(err.Error(), "metaschema") {
		t.Errorf("expected the underlying compile error, got %q", err.Error())
	}
}

// RED: the user-facing message must name the real cause, not "could not find schema".
func TestValidateReportsUnusableSchemaCause(t *testing.T) {
	val := v{
		opts:           Opts{SkipKinds: map[string]struct{}{}, RejectKinds: map[string]struct{}{}},
		schemaDownload: downloadSchema,
		regs:           []registry.Registry{unusableSchemaRegistry()},
	}

	got := val.ValidateResource(resource.Resource{Bytes: []byte("kind: name\napiVersion: v1\n")})

	if got.Status != Error {
		t.Fatalf("expected Error, got %d", got.Status)
	}
	if got.Err == nil {
		t.Fatalf("expected an error")
	}
	if strings.Contains(got.Err.Error(), "could not find schema") {
		t.Errorf("schema was found but could not be compiled; message misreports it as missing: %q", got.Err.Error())
	}
}

// GUARD (green on base too): a genuinely missing schema must stay a missing
// schema - (nil, nil) - so that -ignore-missing-schemas keeps working.
func TestDownloadSchemaMissingSchemaStaysNil(t *testing.T) {
	schema, err := downloadSchema(
		[]registry.Registry{notFoundRegistry(), notFoundRegistry()},
		jsonschema.SchemeURLLoader{}, "name", "v1", "1.30.0")

	if schema != nil || err != nil {
		t.Fatalf("expected (nil, nil) for a missing schema, got (%v, %v)", schema, err)
	}
}

// GUARD (green on base too): -ignore-missing-schemas still skips.
func TestValidateIgnoreMissingSchemasStillSkips(t *testing.T) {
	val := v{
		opts:           Opts{SkipKinds: map[string]struct{}{}, RejectKinds: map[string]struct{}{}, IgnoreMissingSchemas: true},
		schemaDownload: downloadSchema,
		regs:           []registry.Registry{notFoundRegistry()},
	}

	got := val.ValidateResource(resource.Resource{Bytes: []byte("kind: name\napiVersion: v1\n")})
	if got.Status != Skipped {
		t.Fatalf("expected Skipped, got %d (err=%v)", got.Status, got.Err)
	}
}

// GUARD (green on base too): the first usable schema still wins when an earlier
// registry served an unusable one.
func TestDownloadSchemaFallsBackToNextRegistry(t *testing.T) {
	good := newMockRegistry(func() (string, any, error) {
		s, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(`{"type": "object"}`)))
		return "good.json", s, err
	})

	schema, err := downloadSchema(
		[]registry.Registry{unusableSchemaRegistry(), good},
		jsonschema.SchemeURLLoader{}, "name", "v1", "1.30.0")

	if err != nil {
		t.Fatalf("expected the good schema to win, got error %v", err)
	}
	if schema == nil {
		t.Fatalf("expected a schema from the second registry")
	}
}
