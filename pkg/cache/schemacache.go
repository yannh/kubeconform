package cache

import (
	"fmt"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaCache is a cache for downloaded schemas, so each file is only retrieved once
type SchemaCache struct {
	sync.RWMutex
	schemas map[string]*gojsonschema.Schema
}

// New creates a new cache for downloaded schemas
func New() *SchemaCache {
	return &SchemaCache{
		schemas: map[string]*gojsonschema.Schema{},
	}
}

// Key computes a key for a specific JSON schema from its Kind, the resource API Version, and the
// Kubernetes version
func Key(resourceKind, resourceAPIVersion, k8sVersion string) string {
	return fmt.Sprintf("%s-%s-%s", resourceKind, resourceAPIVersion, k8sVersion)
}

// Get retrieves the JSON schema given a resource signature
func (c *SchemaCache) Get(key string) (*gojsonschema.Schema, bool) {
	c.RLock()
	defer c.RUnlock()
	schema, ok := c.schemas[key]
	return schema, ok
}

// Set adds a JSON schema to the schema cache
func (c *SchemaCache) Set(key string, schema *gojsonschema.Schema) {
	c.Lock()
	defer c.Unlock()
	c.schemas[key] = schema
}
