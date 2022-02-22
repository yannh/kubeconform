package registry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPGetter struct {
	httpGet func(string) (*http.Response, error)
}

func newMockHTTPGetter(f func(string) (*http.Response, error)) *mockHTTPGetter {
	return &mockHTTPGetter{
		httpGet: f,
	}
}
func (m mockHTTPGetter) Get(url string) (resp *http.Response, err error) {
	return m.httpGet(url)
}

func TestDownloadSchema(t *testing.T) {
	for _, testCase := range []struct {
		name                                         string
		c                                            httpGetter
		schemaPathTemplate                           string
		strict                                       bool
		resourceKind, resourceAPIVersion, k8sversion string
		expect                                       []byte
		expectErr                                    error
	}{
		{
			"error when downloading",
			newMockHTTPGetter(func(url string) (resp *http.Response, err error) {
				return nil, fmt.Errorf("failed downloading from registry")
			}),
			"http://kubernetesjson.dev",
			true,
			"Deployment",
			"v1",
			"1.18.0",
			nil,
			fmt.Errorf("failed downloading schema at http://kubernetesjson.dev: failed downloading from registry"),
		},
		{
			"getting 404",
			newMockHTTPGetter(func(url string) (resp *http.Response, err error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(strings.NewReader("http response mock body")),
				}, nil
			}),
			"http://kubernetesjson.dev",
			true,
			"Deployment",
			"v1",
			"1.18.0",
			nil,
			fmt.Errorf("no schema found"),
		},
		{
			"getting 503",
			newMockHTTPGetter(func(url string) (resp *http.Response, err error) {
				return &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Body:       ioutil.NopCloser(strings.NewReader("http response mock body")),
				}, nil
			}),
			"http://kubernetesjson.dev",
			true,
			"Deployment",
			"v1",
			"1.18.0",
			nil,
			fmt.Errorf("error while downloading schema at http://kubernetesjson.dev - received HTTP status 503"),
		},
		{
			"200",
			newMockHTTPGetter(func(url string) (resp *http.Response, err error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader("http response mock body")),
				}, nil
			}),
			"http://kubernetesjson.dev",
			true,
			"Deployment",
			"v1",
			"1.18.0",
			[]byte("http response mock body"),
			nil,
		},
	} {
		reg := SchemaRegistry{
			c:                  testCase.c,
			schemaPathTemplate: testCase.schemaPathTemplate,
			strict:             testCase.strict,
		}

		res, err := reg.DownloadSchema(testCase.resourceKind, testCase.resourceAPIVersion, testCase.k8sversion)
		if err == nil || testCase.expectErr == nil {
			if err != testCase.expectErr {
				t.Errorf("during test '%s': expected error, got:\n%s\n%s\n", testCase.name, testCase.expectErr, err)
			}
		} else if err.Error() != testCase.expectErr.Error() {
			t.Errorf("during test '%s': expected error, got:\n%s\n%s\n", testCase.name, testCase.expectErr, err)
		}

		if !bytes.Equal(res, testCase.expect) {
			t.Errorf("during test '%s': expected %s, got %s", testCase.name, testCase.expect, res)
		}
	}

}
