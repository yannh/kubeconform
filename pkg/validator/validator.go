package validator

import (
	"fmt"
	"github.com/yannh/kubeconform/pkg/resource"

	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

type Status int

const (
	_ Status = iota
	Error
	Skipped
	Valid
	Invalid
	Empty
)

// ValidFormat is a type for quickly forcing
// new formats on the gojsonschema loader
type ValidFormat struct{}

// ValidFormat is a type for quickly forcing
// new formats on the gojsonschema loader
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

// Result contains the details of the result of a resource validation
type Result struct {
	Resource resource.Resource
	Err      error
	Status   Status
}

// NewError is a utility function to generate a validation error
func NewError(filename string, err error) Result {
	return Result{
		Resource: resource.Resource{Path: filename},
		Err:      err,
		Status:   Error,
	}
}

// Validate validates a single Kubernetes resource against a Json Schema
func Validate(res resource.Resource, schema *gojsonschema.Schema) Result {
	if schema == nil {
		return Result{Resource: res, Status: Skipped, Err: nil}
	}

	var resource map[string]interface{}
	if err := yaml.Unmarshal(res.Bytes, &resource); err != nil {
		return Result{Resource: res, Status: Error, Err: fmt.Errorf("error unmarshalling resource: %s", err)}
	}
	resourceLoader := gojsonschema.NewGoLoader(resource)

	results, err := schema.Validate(resourceLoader)
	if err != nil {
		// This error can only happen if the Object to validate is poorly formed. There's no hope of saving this one
		return Result{Resource: res, Status: Error, Err: fmt.Errorf("problem validating schema. Check JSON formatting: %s", err)}
	}

	if results.Valid() {
		return Result{Resource: res, Status: Valid}
	}

	msg := ""
	for _, errMsg := range results.Errors() {
		if msg != "" {
			msg += " - "
		}
		msg += errMsg.Description()
	}

	return Result{Resource: res, Status: Invalid, Err: fmt.Errorf("%s", msg)}
}
