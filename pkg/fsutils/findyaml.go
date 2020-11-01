package fsutils

import (
	"os"
	"path/filepath"
	"strings"
)

func isYaml(info os.FileInfo) bool {
	return !info.IsDir() && (strings.HasSuffix(strings.ToLower(info.Name()), ".yaml") || strings.HasSuffix(strings.ToLower(info.Name()), ".yml"))
}

// FindYamlInDir will find yaml files in folder dir, and send their filenames in batches
// of size batchSize to channel fileBatches
func FindYamlInDir(dir string, files chan<- string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isYaml(info) {
			files <- path
		}

		return nil
	})
}
