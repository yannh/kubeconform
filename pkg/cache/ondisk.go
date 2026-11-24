package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
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
func (c *onDisk) Get(key string) (any, error) {
	c.RLock()
	defer c.RUnlock()

	f, err := os.Open(cachePath(c.folder, key))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

// Set adds a JSON schema to the schema cache, replacing any existing entry.
//
// The cache folder may be shared by several kubeconform processes, which the
// mutex does not serialise. The schema is therefore written to a temporary
// file and renamed into place, so that a concurrent reader sees either no
// entry or the complete one, never a partially written file. Rename is atomic
// on POSIX; on Windows it is a best-effort replace.
func (c *onDisk) Set(key string, schema any) error {
	c.Lock()
	defer c.Unlock()

	dest := cachePath(c.folder, key)
	tmp, err := os.CreateTemp(c.folder, filepath.Base(dest)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	_, err = tmp.Write(schema.([]byte))
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		// CreateTemp creates the file 0600; keep the previous 0644 permissions.
		err = os.Chmod(tmpName, 0644)
	}
	if err == nil {
		err = os.Rename(tmpName, dest)
	}
	if err != nil {
		os.Remove(tmpName)
		// On Windows the rename can fail while another process has the entry
		// open. The entry then exists already, which is all the cache needs.
		if _, statErr := os.Stat(dest); statErr == nil {
			return nil
		}
		return err
	}
	return nil
}
