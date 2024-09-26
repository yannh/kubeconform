package registry

import (
	"embed"
)

//go:embed *.json
var content embed.FS

type EmbeddedRegistry struct {
	debug  bool
	strict bool
}

// NewEmbeddedRegistry creates a new "registry", that will serve schemas from embedded resource
func NewEmbeddedRegistry(debug bool, strict bool) *EmbeddedRegistry {
	return &EmbeddedRegistry{
		debug, strict,
	}
}

// DownloadSchema retrieves the schema from a file for the resource
func (r EmbeddedRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (string, []byte, error) {
	var fileName string
	if r.strict {
		fileName = resourceKind + "-strict.json"
	} else {
		fileName = resourceKind + ".json"
	}
	bytes, err := content.ReadFile(fileName)
	if err != nil {
		return resourceKind, nil, nil
	}
	return "embedded:" + resourceKind, bytes, nil
}
