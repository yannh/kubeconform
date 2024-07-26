package registry

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
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

	firstAttempt := true

	// http server to simulate different responses
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is asking to simulate a connection reset or 503
		// then return that status code once and then return a normal response
		if r.URL.Path == "/simulate-reset" || r.URL.Path == "/503" {
			if firstAttempt {
				firstAttempt = false
				if r.URL.Path == "/simulate-reset" {
					if hj, ok := w.(http.Hijacker); ok {
						conn, _, err := hj.Hijack()
						if err != nil {
							fmt.Printf("Hijacking failed: %v\n", err)
							return
						}
						conn.Close() // Close the connection to simulate a reset
					}
				}
				if r.URL.Path == "/503" {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("503 Service Unavailable"))
				}
				return
			} else {
				firstAttempt = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Normal response"))
				return
			}
		}

		if r.URL.Path == "/404" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 Not Found"))
			return
		}

		if r.URL.Path == "/500" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		// Serve a normal response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Normal response"))
	})

	port := fmt.Sprint(rand.Intn(1000) + 9000) // random port
	server := &http.Server{Addr: "127.0.0.1:" + port}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()

	url := fmt.Sprintf("http://localhost:%s", port)

	// Wait for the server to start
	for i := 0; i < 20; i++ {
		fmt.Printf("Trying to connect to server %d ...\n", i)
		_, err := http.Get(url)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	for _, testCase := range []struct {
		name                                         string
		schemaPathTemplate                           string
		strict                                       bool
		resourceKind, resourceAPIVersion, k8sversion string
		expect                                       []byte
		expectErr                                    error
	}{
		{
			"retry connection reset by peer",
			fmt.Sprintf("%s/simulate-reset", url),
			true,
			"Deployment",
			"v1",
			"1.18.0",
			[]byte("Normal response"),
			nil,
		},
		{
			"getting 404",
			fmt.Sprintf("%s/404", url),
			true,
			"Deployment",
			"v1",
			"1.18.0",
			nil,
			fmt.Errorf("could not find schema at %s/404", url),
		},
		{
			"getting 500",
			fmt.Sprintf("%s/500", url),
			true,
			"Deployment",
			"v1",
			"1.18.0",
			nil,
			fmt.Errorf("failed downloading schema at %s/500: Get \"%s/500\": GET %s/500 giving up after 3 attempt(s)", url, url, url),
		},
		{
			"retry 503",
			fmt.Sprintf("%s/503", url),
			true,
			"Deployment",
			"v1",
			"1.18.0",
			[]byte("Normal response"),
			nil,
		},
		{
			"200",
			url,
			true,
			"Deployment",
			"v1",
			"1.18.0",
			[]byte("Normal response"),
			nil,
		},
	} {
		// create a temporary directory for the cache
		tmpDir, err := os.MkdirTemp("", "kubeconform-cache")
		if err != nil {
			t.Errorf("during test '%s': failed to create temp directory: %s", testCase.name, err)
			continue
		}
		defer os.RemoveAll(tmpDir) // clean up the temporary directory

		reg, err := newHTTPRegistry(testCase.schemaPathTemplate, tmpDir, testCase.strict, true, true)
		if err != nil {
			t.Errorf("during test '%s': failed to create registry: %s", testCase.name, err)
			continue
		}

		_, res, err := reg.DownloadSchema(testCase.resourceKind, testCase.resourceAPIVersion, testCase.k8sversion)
		if err == nil || testCase.expectErr == nil {
			if err == nil && testCase.expectErr != nil {
				t.Errorf("during test '%s': expected error\n%s, got nil", testCase.name, testCase.expectErr)
			}
			if err != nil && testCase.expectErr == nil {
				t.Errorf("during test '%s': expected no error, got\n%s\n", testCase.name, err)
			}
		} else if err.Error() != testCase.expectErr.Error() {
			t.Errorf("during test '%s': expected error\n%s, got:\n%s\n", testCase.name, testCase.expectErr, err)
		}

		if !bytes.Equal(res, testCase.expect) {
			t.Errorf("during test '%s': expected '%s', got '%s'", testCase.name, testCase.expect, res)
		}
	}

}
