package cache

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"sync"
)

var mu sync.Mutex
var schemas map[string]*gojsonschema.Schema

func init () {
	schemas = map[string]*gojsonschema.Schema{}
}

func WithCache(downloadSchema func(string, string, string) (*gojsonschema.Schema, error)) (func (string, string, string) (*gojsonschema.Schema, error)) {
	return func(resourceKind, resourceAPIVersion, k8sVersion string) (*gojsonschema.Schema, error) {
		cacheKey := fmt.Sprintf("%s-%s-%s", resourceKind, resourceAPIVersion, k8sVersion)
		mu.Lock()
		cachedSchema, ok := schemas[cacheKey];
		mu.Unlock()
		if ok {
			return cachedSchema, nil
		}

		schema, err := downloadSchema(resourceKind, resourceAPIVersion, k8sVersion)

		mu.Lock()
		schemas[cacheKey] = schema
		mu.Unlock()

		return schema, err
	}
}