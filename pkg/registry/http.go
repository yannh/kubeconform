package registry

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// SchemaRegistry is a file repository (local or remote) that contains JSON schemas for Kubernetes resources
type SchemaRegistry struct {
	schemaPathTemplate string
	strict             bool
	debug              bool
	loader             jsonschema.URLLoader
}

func newHTTPRegistry(schemaPathTemplate string, loader jsonschema.URLLoader, strict bool, debug bool) (*SchemaRegistry, error) {
	return &SchemaRegistry{
		schemaPathTemplate: schemaPathTemplate,
		strict:             strict,
		loader:             loader,
		debug:              debug,
	}, nil
}

// DownloadSchema downloads the schema for a particular resource from an HTTP server
func (r SchemaRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (string, any, error) {
	url, err := schemaPath(r.schemaPathTemplate, resourceKind, resourceAPIVersion, k8sVersion, r.strict)
	if err != nil {
		return "", nil, err
	}

	resp, err := r.loader.Load(url)

	return url, resp, err
}
