package validator

import (
	"github.com/yannh/kubeconform/pkg/registry"
	"testing"

	"github.com/yannh/kubeconform/pkg/resource"
)

type mockRegistry struct {
	SchemaDownloader func() ([]byte, error)
}

func newMockRegistry(f func() ([]byte, error)) *mockRegistry {
	return &mockRegistry{
		SchemaDownloader: f,
	}
}

func (m mockRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	return m.SchemaDownloader()
}

func TestValidate(t *testing.T) {
	for i, testCase := range []struct {
		name                         string
		rawResource, schemaRegistry1 []byte
		schemaRegistry2              []byte
		ignoreMissingSchema          bool
		expect                       Status
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
			Valid,
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
			Invalid,
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
			Invalid,
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
			Error,
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
			Valid,
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
			Valid,
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
			Skipped,
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
			Error,
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
			Skipped,
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
			Error,
		},
	} {
		val := v{
			opts: Opts{
				SkipKinds:            map[string]struct{}{},
				RejectKinds:          map[string]struct{}{},
				IgnoreMissingSchemas: testCase.ignoreMissingSchema,
			},
			schemaCache:    nil,
			schemaDownload: downloadSchema,
			regs: []registry.Registry{
				newMockRegistry(func() ([]byte, error) {
					return testCase.schemaRegistry1, nil
				}),
				newMockRegistry(func() ([]byte, error) {
					return testCase.schemaRegistry2, nil
				}),
			},
		}
		if got := val.ValidateResource(resource.Resource{Bytes: testCase.rawResource}); got.Status != testCase.expect {
			t.Errorf("%d - expected %d, got %d: %s", i, testCase.expect, got.Status, got.Err.Error())
		}
	}
}
