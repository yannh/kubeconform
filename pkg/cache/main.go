package cache

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var schemas map[string][]byte

func init () {
	schemas = map[string][]byte {}
}

func WithCache(downloadSchema func(string, string, string) ([]byte, error)) (func (string, string, string) ([]byte, error)) {
	return func(resourceKind, resourceAPIVersion, k8sVersion string) ([]byte, error) {
		cacheKey := fmt.Sprintf("%s-%s-%s", resourceKind, resourceAPIVersion, k8sVersion)
		mu.Lock()
		cachedSchema, ok := schemas[cacheKey];
		mu.Unlock()
		if ok {
			return cachedSchema, nil
		}

		schema, err := downloadSchema(resourceKind, resourceAPIVersion, k8sVersion)
		if err != nil {
			return schema, err
		}

		mu.Lock()
		schemas[cacheKey] = schema
		mu.Unlock()

		return schema, nil
	}
}