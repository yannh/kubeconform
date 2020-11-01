package resource

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func isYAMLFile(info os.FileInfo) bool {
	return !info.IsDir() && (strings.HasSuffix(strings.ToLower(info.Name()), ".yaml") || strings.HasSuffix(strings.ToLower(info.Name()), ".yml"))
}

func isJSONFile(info os.FileInfo) bool {
	return !info.IsDir() && (strings.HasSuffix(strings.ToLower(info.Name()), ".json"))
}

type DiscoveryError struct {
	Path string
	Err  error
}

func (de DiscoveryError) Error() string {
	return de.Err.Error()
}

func FromFiles(paths ...string) (<-chan Resource, <-chan error) {
	resources := make(chan Resource)
	errors := make(chan error)

	go func() {
		for _, path := range paths {
			// we handle errors in the walk function directly
			// so it should be safe to discard the outer error
			_ = filepath.Walk(path, func(p string, i os.FileInfo, err error) error {
				if err != nil {
					errors <- DiscoveryError{p, err}
					return err
				}

				if !isYAMLFile(i) && !isJSONFile(i) {
					return nil
				}

				f, err := os.Open(p)
				if err != nil {
					errors <- DiscoveryError{p, err}
					return err
				}

				b, err := ioutil.ReadAll(f)
				if err != nil {
					errors <- DiscoveryError{p, err}
					return err
				}

				for _, r := range bytes.Split(b, []byte("---\n")) {
					resources <- Resource{Path: p, Bytes: r}
				}

				return nil
			})
		}

		close(resources)
		close(errors)
	}()

	return resources, errors
}
