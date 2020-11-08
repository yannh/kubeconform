package resource

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

func isIgnored(path string, ignoreFilePatterns []string) (bool, error) {
	for _, p := range ignoreFilePatterns {
		m, err := regexp.MatchString(p, path)
		if err != nil {
			return false, err
		}
		if m {
			return true, nil
		}
	}
	return false, nil
}

func FromFiles(ctx context.Context, ignoreFilePatterns []string, paths ...string) (<-chan Resource, <-chan error) {
	resources := make(chan Resource)
	errors := make(chan error)
	stop := false

	go func() {
		<-ctx.Done()
		stop = true
	}()

	go func() {
		for _, path := range paths {
			// we handle errors in the walk function directly
			// so it should be safe to discard the outer error
			err := filepath.Walk(path, func(p string, i os.FileInfo, err error) error {
				if stop == true {
					return io.EOF
				}
				if err != nil {
					return err
				}

				if !isYAMLFile(i) && !isJSONFile(i) {
					return nil
				}

				ignored, err := isIgnored(path, ignoreFilePatterns)
				if err != nil {
					return err
				}
				if ignored {
					return nil
				}

				f, err := os.Open(p)
				if err != nil {
					return err
				}

				b, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}

				for _, r := range bytes.Split(b, []byte("---\n")) {
					resources <- Resource{Path: p, Bytes: r}
				}

				return nil
			})

			if err != nil && err != io.EOF {
				errors <- DiscoveryError{path, err}
			}
		}

		close(resources)
		close(errors)
	}()

	return resources, errors
}
