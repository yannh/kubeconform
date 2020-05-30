package registry

type Manifest struct {
	Kind, Version string
}

type Registry interface {
	DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error)
}
