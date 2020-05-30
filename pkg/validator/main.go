package validator

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

type InvalidResourceError struct { err string }
func (r InvalidResourceError) Error() string{
	return r.err
}

// ValidFormat is a type for quickly forcing
// new formats on the gojsonschema loader
type ValidFormat struct{}
func (f ValidFormat) IsFormat(input interface{}) bool {
	return true
}

// From kubeval - let's see if absolutely necessary
// func init () {
// 	gojsonschema.FormatCheckers.Add("int64", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("byte", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("int32", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("int-or-string", ValidFormat{})
// }

func Validate(rawResource []byte, rawSchema []byte) error {
	schemaLoader := gojsonschema.NewBytesLoader(rawSchema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return err
	}

	var resource map[string]interface{}
	if err = yaml.Unmarshal(rawResource, &resource); err != nil {
		return fmt.Errorf("error unmarshalling resource: %s", err)
	}
	resourceLoader := gojsonschema.NewGoLoader(resource)

	results, err := schema.Validate(resourceLoader)
	if err != nil {
		// This error can only happen if the Object to validate is poorly formed. There's no hope of saving this one
		return  fmt.Errorf("problem validating schema. Check JSON formatting: %s", err)
	}

	if !results.Valid() {
		return InvalidResourceError{err: fmt.Sprintf("resource does not conform to schema: %+v", results) }
	}

	return nil
}

