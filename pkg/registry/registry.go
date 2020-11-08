package registry

import (
	"bytes"
	"strings"
	"text/template"
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

func schemaPath(tpl, resourceKind, resourceAPIVersion, k8sVersion string, strict bool) (string, error) {
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

	tmpl, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", err
	}

	tplData := struct {
		NormalizedVersion string
		StrictSuffix      string
		ResourceKind      string
		KindSuffix        string
	}{
		normalisedVersion,
		strictSuffix,
		strings.ToLower(resourceKind),
		kindSuffix,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tplData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func New(schemaLocation string, strict bool, skipTLS bool) Registry {
	if !strings.HasSuffix(schemaLocation, "json") { // If we dont specify a full templated path, we assume the paths of kubernetesjsonschema.dev
		schemaLocation += "/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json"
	}

	if strings.HasPrefix(schemaLocation, "http") {
		return newHTTPRegistry(schemaLocation, strict, skipTLS)
	} else {
		return newLocalRegistry(schemaLocation, strict)
	}
}
