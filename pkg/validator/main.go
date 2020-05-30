package validator

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

// ValidFormat is a type for quickly forcing
// new formats on the gojsonschema loader
type ValidFormat struct{}
func (f ValidFormat) IsFormat(input interface{}) bool {
	return true
}

func init () {
	gojsonschema.FormatCheckers.Add("int64", ValidFormat{})
	gojsonschema.FormatCheckers.Add("byte", ValidFormat{})
	gojsonschema.FormatCheckers.Add("int32", ValidFormat{})
	gojsonschema.FormatCheckers.Add("int-or-string", ValidFormat{})
}

func Validate(resource interface{}, rawSchema []byte) error {
	schemaLoader := gojsonschema.NewBytesLoader(rawSchema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return err
	}

	documentLoader := gojsonschema.NewGoLoader(resource)
	results, err := schema.Validate(documentLoader)
	if err != nil {
		// This error can only happen if the Object to validate is poorly formed. There's no hope of saving this one
		return  fmt.Errorf("problem validating schema. Check JSON formatting: %s", err)
	}

	if !results.Valid() {
		return fmt.Errorf("resource does not conform to schema")
	}

	return nil
}

