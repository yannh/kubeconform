package loader

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type mockCache struct {
	data map[string]any
}

func (m *mockCache) Get(key string) (any, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, errors.New("cache miss")
}

func (m *mockCache) Set(key string, value any) error {
	m.data[key] = value
	return nil
}

// Test basic functionality of HTTPURLLoader
func TestHTTPURLLoader_Load(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		cacheEnabled   bool
		expectError    bool
		expectCacheHit bool
	}{
		{
			name:           "successful load",
			mockResponse:   `{"type": "object"}`,
			mockStatusCode: http.StatusOK,
			cacheEnabled:   false,
			expectError:    false,
		},
		{
			name:           "not found error",
			mockResponse:   "",
			mockStatusCode: http.StatusNotFound,
			cacheEnabled:   false,
			expectError:    true,
		},
		{
			name:           "server error",
			mockResponse:   "",
			mockStatusCode: http.StatusInternalServerError,
			cacheEnabled:   false,
			expectError:    true,
		},
		{
			name:           "cache hit",
			mockResponse:   `{"type": "object"}`,
			mockStatusCode: http.StatusOK,
			cacheEnabled:   true,
			expectError:    false,
			expectCacheHit: true,
		},
		{
			name:           "Partial response from server",
			mockResponse:   `{"type": "objec`,
			mockStatusCode: http.StatusOK,
			cacheEnabled:   false,
			expectError:    true,
			expectCacheHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create HTTPURLLoader
			loader := &HTTPURLLoader{
				client: *server.Client(),
				cache:  nil,
			}

			if tt.cacheEnabled {
				loader.cache = &mockCache{data: map[string]any{}}
				if tt.expectCacheHit {
					loader.cache.Set(server.URL, []byte(tt.mockResponse))
				}
			}

			// Call Load and handle errors
			res, err := loader.Load(server.URL)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if res == nil {
					t.Errorf("expected non-nil result, got nil")
				}
			}

		})
	}
}

// Test basic functionality of HTTPURLLoader
func TestHTTPURLLoader_Load_Retries(t *testing.T) {

	tests := []struct {
		name                string
		url                 string
		expectError         bool
		expectCallCount     int
		consecutiveFailures int
	}{
		{
			name:                "retries on 503",
			url:                 "/503",
			expectError:         false,
			expectCallCount:     2,
			consecutiveFailures: 2,
		},
		{
			name:                "fails when hitting max retries",
			url:                 "/503",
			expectError:         true,
			expectCallCount:     3,
			consecutiveFailures: 5,
		},
		{
			name:                "retry on connection reset",
			url:                 "/simulate-reset",
			expectError:         false,
			expectCallCount:     2,
			consecutiveFailures: 1,
		},
		{
			name:                "retry on connection reset",
			url:                 "/simulate-reset",
			expectError:         true,
			expectCallCount:     3,
			consecutiveFailures: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ccMutex := &sync.Mutex{}
			callCounts := map[string]int{}
			// Mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ccMutex.Lock()
				callCounts[r.URL.Path]++
				callCount := callCounts[r.URL.Path]
				ccMutex.Unlock()

				switch r.URL.Path {
				case "/simulate-reset":
					if callCount <= tt.consecutiveFailures {
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

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"type": "object"}`))

				case "/503":
					s := http.StatusServiceUnavailable
					if callCount >= tt.consecutiveFailures {
						s = http.StatusOK
					}
					w.WriteHeader(s)
					w.Write([]byte(`{"type": "object"}`))
				}
			}))
			defer server.Close()

			// Create HTTPURLLoader
			loader, _ := NewHTTPURLLoader(false, nil, "")

			fullurl := server.URL + tt.url
			// Call Load and handle errors
			_, err := loader.Load(fullurl)
			if tt.expectError && err == nil {
				t.Error("expected error, got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			ccMutex.Lock()
			if callCounts[tt.url] != tt.expectCallCount {
				t.Errorf("expected %d calls, got: %d", tt.expectCallCount, callCounts[tt.url])
			}
			ccMutex.Unlock()
		})
	}
}
