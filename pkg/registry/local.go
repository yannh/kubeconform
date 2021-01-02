package registry

import (
	"fmt"
	"io/ioutil"
	"os"
)

type LocalRegistry struct {
	pathTemplate string
	strict       bool
}

// NewLocalSchemas creates a new "registry", that will serve schemas from files, given a list of schema filenames
func newLocalRegistry(pathTemplate string, strict bool) (*LocalRegistry, error) {
	return &LocalRegistry{
		pathTemplate,
		strict,
	}, nil
}

// DownloadSchema retrieves the schema from a file for the resource
func (r LocalRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	schemaFile, err := schemaPath(r.pathTemplate, resourceKind, resourceAPIVersion, k8sVersion, r.strict)
	if err != nil {
		return []byte{}, nil
	}
	f, err := os.Open(schemaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, newNotFoundError(fmt.Errorf("no schema found"))
		}
		return nil, fmt.Errorf("failed to open schema %s", schemaFile)
	}

	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return content, nil
}
