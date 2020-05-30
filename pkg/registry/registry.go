package registry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Manifest struct {
	Kind, Version string
}

type Registry interface {
	DownloadSchema(kind, apiVersion string) error
}

type KubernetesRegistry struct {
	baseURL    string
	strict     bool
}

func NewKubernetesRegistry(strict bool) *KubernetesRegistry {
	return &KubernetesRegistry{
		baseURL:    "https://kubernetesjsonschema.dev",
		strict:     strict,
	}
}


func (r KubernetesRegistry) schemaURL(resourceKind, resourceAPIVersion, k8sVersion string) string {
	normalisedVersion := k8sVersion
	if normalisedVersion != "master" {
		normalisedVersion = "v" + normalisedVersion
	}

	strictSuffix := ""
	if r.strict {
		strictSuffix = "-strict"
	}

	groupParts := strings.Split(resourceAPIVersion, "/")
	versionParts := strings.Split(groupParts[0], ".")

	kindSuffix := "-" + strings.ToLower(versionParts[0])
	if len(groupParts) > 1 {
		kindSuffix += "-" + strings.ToLower(groupParts[1])
	}

	return fmt.Sprintf("%s/%s-standalone%s/%s%s.json", r.baseURL, normalisedVersion, strictSuffix, strings.ToLower(resourceKind), kindSuffix)
}

func (r KubernetesRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	url := r.schemaURL(resourceKind, resourceAPIVersion, k8sVersion)

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}
	defer resp.Body.Close()

	body, err :=  ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}

	fmt.Printf("downloaded %s\n", url)

	return body, nil
}

