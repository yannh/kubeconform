package validator

import (
	"testing"

	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"

	"github.com/xeipuuv/gojsonschema"
)

func TestValidate(t *testing.T) {
	for i, testCase := range []struct {
		name                string
		rawResource, schema []byte
		expect              Status
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
			Error,
		},
	} {
		val := v{
			opts: Opts{
				SkipKinds:   map[string]struct{}{},
				RejectKinds: map[string]struct{}{},
			},
			schemaCache: nil,
			schemaDownload: func(_ []registry.Registry, _, _, _ string) (*gojsonschema.Schema, error) {
				schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(testCase.schema))
				if err != nil {
					t.Errorf("failed parsing test schema")
				}
				return schema, nil
			},
			regs: nil,
		}
		if got := val.ValidateResource(resource.Resource{Bytes: testCase.rawResource}); got.Status != testCase.expect {
			t.Errorf("%d - expected %d, got %d: %s", i, testCase.expect, got.Status, got.Err.Error())
		}
	}
}
