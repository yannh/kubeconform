package loader

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/yannh/kubeconform/pkg/cache"
	"io"
	"net/http"
	"time"
)

type HTTPURLLoader struct {
	client http.Client
	cache  cache.Cache
}

func (l *HTTPURLLoader) Load(url string) (any, error) {
	if l.cache != nil {
		if cached, err := l.cache.Get(url); err == nil {
			return jsonschema.UnmarshalJSON(bytes.NewReader(cached.([]byte)))
		}
	}

	resp, err := l.client.Get(url)
	if err != nil {
		msg := fmt.Sprintf("failed downloading schema at %s: %s", url, err)
		return nil, errors.New(msg)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		msg := fmt.Sprintf("could not find schema at %s", url)
		return nil, NewNotFoundError(errors.New(msg))
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("error while downloading schema at %s - received HTTP status %d", url, resp.StatusCode)
		return nil, fmt.Errorf("%s", msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("failed parsing schema from %s: %s", url, err)
		return nil, errors.New(msg)
	}

	if l.cache != nil {
		if err = l.cache.Set(url, body); err != nil {
			return nil, fmt.Errorf("failed to write cache to disk: %s", err)
		}
	}

	s, err := jsonschema.UnmarshalJSON(bytes.NewReader(body))
	if err != nil {
		return nil, NewNonJSONResponseError(err)
	}

	return s, nil
}

func NewHTTPURLLoader(skipTLS bool, cache cache.Cache) (*HTTPURLLoader, error) {
	transport := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
		Proxy:              http.ProxyFromEnvironment,
	}

	if skipTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// retriable http client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 2
	retryClient.HTTPClient = &http.Client{Transport: transport}
	retryClient.Logger = nil

	httpLoader := HTTPURLLoader{client: *retryClient.StandardClient(), cache: cache}
	return &httpLoader, nil
}
