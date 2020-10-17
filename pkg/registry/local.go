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

type fileNotFoundError struct {
	err         error
	isRetryable bool
}

func newFileNotFoundError(err error, isRetryable bool) *fileNotFoundError {
	return &fileNotFoundError{err, isRetryable}
}
func (e *fileNotFoundError) IsRetryable() bool { return e.isRetryable }
func (e *fileNotFoundError) Error() string     { return e.err.Error() }

// NewLocalSchemas creates a new "registry", that will serve schemas from files, given a list of schema filenames
func NewLocalRegistry(pathTemplate string, strict bool) *LocalRegistry {
	return &LocalRegistry{
		pathTemplate,
		strict,
	}
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
			return nil, newFileNotFoundError(fmt.Errorf("no schema found"), false)
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
