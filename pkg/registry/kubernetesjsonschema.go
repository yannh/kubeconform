package registry

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type KubernetesRegistry struct {
	schemaPathTemplate string
	strict             bool
}

type downloadError struct {
	err         error
	isRetryable bool
}

func newDownloadError(err error, isRetryable bool) *downloadError {
	return &downloadError{err, isRetryable}
}
func (e *downloadError) IsRetryable() bool { return e.isRetryable }
func (e *downloadError) Error() string     { return e.err.Error() }

func newHTTPRegistry(schemaPathTemplate string, strict bool) *KubernetesRegistry {
	return &KubernetesRegistry{
		schemaPathTemplate: schemaPathTemplate,
		strict:             strict,
	}
}

func (r KubernetesRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	url, err := schemaPath(r.schemaPathTemplate, resourceKind, resourceAPIVersion, k8sVersion, r.strict)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, newDownloadError(fmt.Errorf("no schema found"), false)
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
