package registry

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPGetter struct {
	callNumber int
	httpGet    func(mockHTTPGetter, string) (*http.Response, error)
}

func newMockHTTPGetter(f func(mockHTTPGetter, string) (*http.Response, error)) *mockHTTPGetter {
	return &mockHTTPGetter{
		callNumber: 0,
		httpGet:    f,
	}
}
func (m *mockHTTPGetter) Get(url string) (resp *http.Response, err error) {
	m.callNumber = m.callNumber + 1
	return m.httpGet(*m, url)
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
			"retry connection reset by peer",
			newMockHTTPGetter(func(c mockHTTPGetter, url string) (resp *http.Response, err error) {
				if c.callNumber == 1 {
					return nil, &net.OpError{Err: errors.New("connection reset by peer")}
				} else {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("http response mock body")),
					}, nil
				}
			}),
			"http://kubernetesjson.dev",
			true,
			"Deployment",
			"v1",
			"1.18.0",
			[]byte("http response mock body"),
			nil,
		},
		{
			"getting 404",
			newMockHTTPGetter(func(c mockHTTPGetter, url string) (resp *http.Response, err error) {
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
			"getting 500",
			newMockHTTPGetter(func(c mockHTTPGetter, url string) (resp *http.Response, err error) {
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
			"retry 503",
			newMockHTTPGetter(func(c mockHTTPGetter, url string) (resp *http.Response, err error) {
				if c.callNumber == 1 {
					return &http.Response{
						StatusCode: http.StatusServiceUnavailable,
						Body:       io.NopCloser(strings.NewReader("503 http response mock body")),
					}, nil
				} else {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("http response mock body")),
					}, nil
				}
			}),
			"http://kubernetesjson.dev",
			true,
			"Deployment",
			"v1",
			"1.18.0",
			[]byte("http response mock body"),
			nil,
		},
		{
			"200",
			newMockHTTPGetter(func(c mockHTTPGetter, url string) (resp *http.Response, err error) {
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
