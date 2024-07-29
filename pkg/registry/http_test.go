package registry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPDoer struct {
	httpDo func(*http.Request) (*http.Response, error)
}

func newMockHTTPDoer(f func(*http.Request) (*http.Response, error)) *mockHTTPDoer {
	return &mockHTTPDoer{
		httpDo: f,
	}
}
func (m mockHTTPDoer) Do(req *http.Request) (resp *http.Response, err error) {
	return m.httpDo(req)
}

func TestDownloadSchema(t *testing.T) {
	for _, testCase := range []struct {
		name                                         string
		c                                            httpDoer
		schemaPathTemplate                           string
		strict                                       bool
		resourceKind, resourceAPIVersion, k8sversion string
		expect                                       []byte
		expectErr                                    error
	}{
		{
			"error when downloading",
			newMockHTTPDoer(func(req *http.Request) (resp *http.Response, err error) {
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
			newMockHTTPDoer(func(req *http.Request) (resp *http.Response, err error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("http response mock body")),
				}, nil
			}),
			"http://kubernetesjson.dev",
			true,
			"Deployment",
			"v1",
			"1.18.0",
			nil,
			fmt.Errorf("could not find schema at http://kubernetesjson.dev"),
		},
		{
			"getting 503",
			newMockHTTPDoer(func(req *http.Request) (resp *http.Response, err error) {
				return &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Body:       io.NopCloser(strings.NewReader("http response mock body")),
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
			newMockHTTPDoer(func(req *http.Request) (resp *http.Response, err error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("http response mock body")),
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

		_, res, err := reg.DownloadSchema(testCase.resourceKind, testCase.resourceAPIVersion, testCase.k8sversion)
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
