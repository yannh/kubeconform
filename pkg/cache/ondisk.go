package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"sync"
)

type onDisk struct {
	sync.RWMutex
	folder string
}

// New creates a new cache for downloaded schemas
func NewOnDiskCache(cache string) Cache {
	return &onDisk{
		folder: cache,
	}
}

func cachePath(folder, key string) string {
	hash := sha256.Sum256([]byte(key))
	return path.Join(folder, hex.EncodeToString(hash[:]))
}

// Get retrieves the JSON schema given a resource signature
func (c *onDisk) Get(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	f, err := os.Open(cachePath(c.folder, key))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

// Set adds a JSON schema to the schema cache
func (c *onDisk) Set(key string, schema interface{}) error {
	c.Lock()
	defer c.Unlock()
	return os.WriteFile(cachePath(c.folder, key), schema.([]byte), 0644)
}
