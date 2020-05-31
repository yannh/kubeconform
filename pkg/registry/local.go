package registry

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

type LocalSchemas struct {
	schemas map[string]*gojsonschema.Schema
}

func NewLocalSchemas(schemaFiles []string) (*LocalSchemas, error) {
	schemas := &LocalSchemas{
		schemas: map[string]*gojsonschema.Schema{},
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
		err = yaml.Unmarshal(content, &parsedSchema)
		if err != nil {
			return nil, fmt.Errorf("failed parsing schema %s", schemaFile)
		}

		if strings.HasSuffix(schemaFile, ".yml") || strings.HasSuffix(schemaFile, ".yaml") {
			asJSON, err := yaml.YAMLToJSON(content)
			if err != nil {
				return nil, fmt.Errorf("error converting manifest %s to JSON: %s", schemaFile, err)
			}

			schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(asJSON))
			if err != nil {
				return nil, err
			}
			schemas.schemas[parsedSchema.Spec.Names.Kind] = schema
		}
	}

	return schemas, nil
}

func (r LocalSchemas) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (*gojsonschema.Schema, error) {
	schema, ok := r.schemas[resourceKind]
	if !ok {
		return nil, fmt.Errorf("no local schema for Kind %s", resourceKind)
	}
	return schema, nil
}
