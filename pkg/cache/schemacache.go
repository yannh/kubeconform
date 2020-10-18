package cache

import (
	"fmt"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type SchemaCache struct {
	sync.RWMutex
	schemas map[string]*gojsonschema.Schema
}

func New() *SchemaCache {
	return &SchemaCache{
		schemas: map[string]*gojsonschema.Schema{},
	}
}

func Key(resourceKind, resourceAPIVersion, k8sVersion string) string {
	return fmt.Sprintf("%s-%s-%s", resourceKind, resourceAPIVersion, k8sVersion)
}

func (c *SchemaCache) Get(key string) (*gojsonschema.Schema, bool) {
	c.RLock()
	defer c.RUnlock()
	schema, ok := c.schemas[key]
	return schema, ok
}

func (c *SchemaCache) Set(key string, schema *gojsonschema.Schema) {
	c.Lock()
	defer c.Unlock()
	c.schemas[key] = schema
}
