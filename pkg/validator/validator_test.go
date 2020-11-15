package validator

import (
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"testing"

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
Kind: name
firstName: foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "Kind": {
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
Kind: name
firstName: foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "Kind": {
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
Kind: name
firstName: foo
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "Kind": {
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
Kind: name
firstName foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "Kind": {
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
				SkipKinds:   map[string]bool{},
				RejectKinds: map[string]bool{},
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
			t.Errorf("%d - expected %d, got %d", i, testCase.expect, got.Status)
		}
	}
}
