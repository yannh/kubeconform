package fsutils

import (
	"os"
	"path/filepath"
	"strings"
)

// FindYamlInDir will find yaml files in folder dir, and send their filenames in batches
// of size batchSize to channel fileBatches
func FindYamlInDir(dir string, fileBatches chan<- string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			fileBatches <- path
		}
		return nil
	})

	return err
}
