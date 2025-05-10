package cache

import (
	"fmt"
	"sync"
)

// SchemaCache is a cache for downloaded schemas, so each file is only retrieved once
// It is different from pkg/registry/http_cache.go in that:
//   - This cache caches the parsed Schemas
type inMemory struct {
	sync.RWMutex
	schemas map[string]interface{}
}

// New creates a new cache for downloaded schemas
func NewInMemoryCache() Cache {
	return &inMemory{
		schemas: map[string]interface{}{},
	}
}

// Get retrieves the JSON schema given a resource signature
func (c *inMemory) Get(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()
	schema, ok := c.schemas[key]

	if !ok {
		return nil, fmt.Errorf("schema not found in in-memory cache")
	}

	return schema, nil
}

// Set adds a JSON schema to the schema cache
func (c *inMemory) Set(key string, schema interface{}) error {
	c.Lock()
	defer c.Unlock()
	c.schemas[key] = schema

	return nil
}
