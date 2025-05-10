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
	schemas map[string]any
}

// New creates a new cache for downloaded schemas
func NewInMemoryCache() Cache {
	return &inMemory{
		schemas: make(map[string]any),
	}
}

// Get retrieves the JSON schema given a resource signature
func (c *inMemory) Get(key string) (any, error) {
	c.RLock()
	defer c.RUnlock()
	schema, ok := c.schemas[key]

	if !ok {
		return nil, fmt.Errorf("schema not found in in-memory cache")
	}

	return schema, nil
}

// Set adds a JSON schema to the schema cache
func (c *inMemory) Set(key string, schema any) error {
	c.Lock()
	defer c.Unlock()
	c.schemas[key] = schema

	return nil
}
