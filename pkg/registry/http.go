package registry

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/yannh/kubeconform/pkg/cache"
)

type httpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

// SchemaRegistry is a file repository (local or remote) that contains JSON schemas for Kubernetes resources
type SchemaRegistry struct {
	c                  httpGetter
	schemaPathTemplate string
	cache              cache.Cache
	strict             bool
}

func newHTTPRegistry(schemaPathTemplate string, cacheFolder string, strict bool, skipTLS bool) *SchemaRegistry {
	reghttp := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}

	if skipTLS {
		reghttp.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	var filecache cache.Cache = nil
	if cacheFolder != "" {
		filecache = cache.NewOnDiskCache(cacheFolder)
	}

	return &SchemaRegistry{
		c:                  &http.Client{Transport: reghttp},
		schemaPathTemplate: schemaPathTemplate,
		cache:              filecache,
		strict:             strict,
	}
}

// DownloadSchema downloads the schema for a particular resource from an HTTP server
func (r SchemaRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
	url, err := schemaPath(r.schemaPathTemplate, resourceKind, resourceAPIVersion, k8sVersion, r.strict)
	if err != nil {
		return nil, err
	}

	if r.cache != nil {
		if b, err := r.cache.Get(resourceKind, resourceAPIVersion, k8sVersion); err == nil {
			return b.([]byte), nil
		}
	}

	resp, err := r.c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, newNotFoundError(fmt.Errorf("no schema found"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while downloading schema - received HTTP status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}

	if r.cache != nil {
		if err := r.cache.Set(resourceKind, resourceAPIVersion, k8sVersion, body); err != nil {
			return nil, fmt.Errorf("failed writing schema to cache: %s", err)
		}
	}

	return body, nil
}
