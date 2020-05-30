package registry

import "github.com/xeipuuv/gojsonschema"

type Manifest struct {
	Kind, Version string
}

type Registry interface {
	DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (*gojsonschema.Schema, error)
}

type Retryable interface {
	IsRetryable() bool
}