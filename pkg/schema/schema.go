package schema

import (
	"errors"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type multiError []error

func (m multiError) Error() string {
	var r []string
	for _, e := range ([]error)(m) {
		r = append(r, e.Error())
	}
	return strings.Join(r, "\n")
}

// Schema TODO
type Schema gojsonschema.Schema

// Validate TODO
func (s *Schema) Validate(resource map[string]interface{}) error {
	results, err := (*gojsonschema.Schema)(s).Validate(gojsonschema.NewGoLoader(resource))
	if err != nil {
		return err
	}

	if results.Valid() {
		return nil
	}

	var errs []error
	for _, e := range results.Errors() {
		errs = append(errs, errors.New(e.Description()))
	}

	return multiError(errs)
}
