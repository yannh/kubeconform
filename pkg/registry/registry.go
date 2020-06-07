package registry

import (
	"fmt"
	"strings"
)

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

func schemaPath(resourceKind, resourceAPIVersion, k8sVersion string, strict bool) string {
	normalisedVersion := k8sVersion
	if normalisedVersion != "master" {
		normalisedVersion = "v" + normalisedVersion
	}

	strictSuffix := ""
	if strict {
		strictSuffix = "-strict"
	}

	groupParts := strings.Split(resourceAPIVersion, "/")
	versionParts := strings.Split(groupParts[0], ".")

	kindSuffix := "-" + strings.ToLower(versionParts[0])
	if len(groupParts) > 1 {
		kindSuffix += "-" + strings.ToLower(groupParts[1])
	}

	return fmt.Sprintf("%s-standalone%s/%s%s.json", normalisedVersion, strictSuffix, strings.ToLower(resourceKind), kindSuffix)
}
