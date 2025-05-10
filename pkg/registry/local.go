package registry

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type LocalRegistry struct {
	pathTemplate string
	strict       bool
	debug        bool
	loader       jsonschema.URLLoader
}

// NewLocalSchemas creates a new "registry", that will serve schemas from files, given a list of schema filenames
func newLocalRegistry(pathTemplate string, loader jsonschema.URLLoader, strict bool, debug bool) (*LocalRegistry, error) {
	return &LocalRegistry{
		pathTemplate,
		strict,
		debug,
		loader,
	}, nil
}

// DownloadSchema retrieves the schema from a file for the resource
func (r LocalRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (string, any, error) {
	schemaFile, err := schemaPath(r.pathTemplate, resourceKind, resourceAPIVersion, k8sVersion, r.strict)
	if err != nil {
		return schemaFile, []byte{}, nil
	}

	s, err := r.loader.Load(schemaFile)
	return schemaFile, s, err
}
