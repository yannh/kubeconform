package registry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type LocalRegistry struct {
	folder string
	strict bool
}

// NewLocalSchemas creates a new "registry", that will serve schemas from files, given a list of schema filenames
func NewLocalRegistry(folder string, strict bool) (*LocalRegistry, error) {
	return &LocalRegistry{
		folder,
		strict,
	}, nil
}

// DownloadSchema retrieves the schema from a file for the resource
func (r LocalRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	schemaFile := path.Join(r.folder, schemaPath(resourceKind, resourceAPIVersion, k8sVersion, r.strict))

	f, err := os.Open(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open schema %s", schemaFile)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return content, nil
}
