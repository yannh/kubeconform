package validator

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"testing"
)

func TestValidate(t *testing.T) {

	for i, testCase := range []struct {
		name                string
		rawResource, schema []byte
		expect              error
	}{
		{
			"valid resource",
			[]byte(`
firstName: foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
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
		},
		{
			"invalid resource",
			[]byte(`
firstName: foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
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
			fmt.Errorf("Invalid type. Expected: number, given: string"),
		},
		{
			"missing required field",
			[]byte(`
firstName: foo
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
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
			fmt.Errorf("lastName is required"),
		},
		{
			"resource has invalid yaml",
			[]byte(`
firstName foo
lastName: bar
`),
			[]byte(`{
  "title": "Example Schema",
  "type": "object",
  "properties": {
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
			fmt.Errorf("error unmarshalling resource: error converting YAML to JSON: yaml: line 3: mapping values are not allowed in this context"),
		},
	} {
		schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(testCase.schema))
		if err != nil {
			t.Errorf("failed parsing test schema")
		}
		if got := Validate(testCase.rawResource, schema); ((got == nil) != (testCase.expect == nil)) || (got != nil && (got.Error() != testCase.expect.Error())) {
			t.Errorf("%d - expected %s, got %s", i, testCase.expect, got)
		}
	}
}
