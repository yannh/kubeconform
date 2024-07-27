package registry

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
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
	callCounts := map[string]int{}

	// http server to simulate different responses
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var s int
		callCounts[r.URL.Path]++
		callCount := callCounts[r.URL.Path]

		switch r.URL.Path {
		case "/404":
			s = http.StatusNotFound
		case "/500":
			s = http.StatusInternalServerError
		case "/503":
			if callCount < 2 {
				s = http.StatusServiceUnavailable
			} else {
				s = http.StatusOK // Should succeed on 3rd try
			}

		case "/simulate-reset":
			if callCount < 2 {
				if hj, ok := w.(http.Hijacker); ok {
					conn, _, err := hj.Hijack()
					if err != nil {
						fmt.Printf("Hijacking failed: %v\n", err)
						return
					}
					conn.Close() // Close the connection to simulate a reset
				}
				return
			}
			s = http.StatusOK // Should succeed on third try

		default:
			s = http.StatusOK
		}

		w.WriteHeader(s)
		w.Write([]byte(http.StatusText(s)))
	})

	port := fmt.Sprint(rand.Intn(1000) + 9000) // random port
	server := &http.Server{Addr: "127.0.0.1:" + port}
	url := fmt.Sprintf("http://localhost:%s", port)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()
	defer server.Shutdown(nil)

	// Wait for the server to start
	for i := 0; i < 20; i++ {
		if _, err := http.Get(url); err == nil {
			break
		}

		if i == 19 {
			t.Error("http server did not start")
			return
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
			[]byte(http.StatusText(http.StatusOK)),
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
			[]byte(http.StatusText(http.StatusOK)),
			nil,
		},
		{
			"200",
			url,
			true,
			"Deployment",
			"v1",
			"1.18.0",
			[]byte(http.StatusText(http.StatusOK)),
			nil,
		},
	} {
		callCounts = map[string]int{} // Reinitialise counters

		reg, err := newHTTPRegistry(testCase.schemaPathTemplate, "", testCase.strict, true, true)
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
