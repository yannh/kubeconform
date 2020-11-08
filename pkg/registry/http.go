package registry

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type KubernetesRegistry struct {
	c                  *http.Client
	schemaPathTemplate string
	strict             bool
}

type NotFoundError struct {
	err error
}

func newNetFoundError(err error) *NotFoundError {
	return &NotFoundError{err}
}
func (e *NotFoundError) Error() string { return e.err.Error() }

func newHTTPRegistry(schemaPathTemplate string, strict bool, skipTLS bool) *KubernetesRegistry {
	reghttp := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}

	if skipTLS {
		reghttp.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &KubernetesRegistry{
		c:                  &http.Client{Transport: reghttp},
		schemaPathTemplate: schemaPathTemplate,
		strict:             strict,
	}
}

func (r KubernetesRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	url, err := schemaPath(r.schemaPathTemplate, resourceKind, resourceAPIVersion, k8sVersion, r.strict)
	if err != nil {
		return nil, err
	}

	resp, err := r.c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, newNetFoundError(fmt.Errorf("no schema found"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while downloading schema - received HTTP status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}

	return body, nil
}
