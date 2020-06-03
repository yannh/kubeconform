package registry

type Manifest struct {
	Kind, Version string
}

// Registry is an interface that should be implemented by any source of Kubernetes schemas
type Registry interface {
	DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error)
}

// Retryable indicates whether an error is a temporary or a permanent failure
type Retryable interface {
	IsRetryable() bool
}
