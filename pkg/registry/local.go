package registry

import (
	"fmt"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

type LocalSchemas struct {
	schemas map[string]string
}

// NewLocalSchemas creates a new "registry", that will serve schemas from files, given a list of schema filenames
func NewLocalSchemas(schemaFiles []string) (*LocalSchemas, error) {
	schemas := &LocalSchemas{
		schemas: map[string]string{},
	}

	for _, schemaFile := range schemaFiles {
		f, err := os.Open(schemaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open schema %s", schemaFile)
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read schema %s", schemaFile)
		}

		var parsedSchema struct {
			Spec struct {
				Names struct {
					Kind string `json:"Kind"`
				} `json:"Names"`
			} `json:"Spec"`
		}
		err = yaml.Unmarshal(content, &parsedSchema) // Index Schemas by kind
		if err != nil {
			return nil, fmt.Errorf("failed parsing schema %s", schemaFile)
		}

		schemas.schemas[parsedSchema.Spec.Names.Kind] = schemaFile
	}

	return schemas, nil
}

// DownloadSchema retrieves the schema from a file for the resource
func (r LocalSchemas) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	schemaFile, ok := r.schemas[resourceKind]
	if !ok {
		return nil, fmt.Errorf("no local schema for Kind %s", resourceKind)
	}

	f, err := os.Open(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open schema %s", schemaFile)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	asJSON := content
	if strings.HasSuffix(schemaFile, ".yml") || strings.HasSuffix(schemaFile, ".yaml") {
		asJSON, err = yaml.YAMLToJSON(content)
		if err != nil {
			return nil, fmt.Errorf("error converting manifest %s to JSON: %s", schemaFile, err)
		}
	}
	return asJSON, nil
}
