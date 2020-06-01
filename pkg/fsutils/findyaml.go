package fsutils

import (
	"os"
	"path/filepath"
	"strings"
)

func FindYamlInDir(dir string, fileBatches chan<- []string, batchSize int) error {
	files := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			files = append(files, path)
			if len(files) > batchSize {
				fileBatches <- files
				files = []string{}
			}
		}
		return nil
	})

	if len(files) > 0 {
		fileBatches <- files
	}

	return err
}
