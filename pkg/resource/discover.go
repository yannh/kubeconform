package resource

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const batchSize = 10

func isYAMLFile(info os.FileInfo) bool {
	return !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml"))
}

func Discover(paths ...string) <-chan []Resource {
	res := make(chan []Resource)

	go func() {
		for _, path := range paths {
			var batch []Resource

			// we handle errors in the walk function directly
			// so it should be safe to discard the outer error
			_ = filepath.Walk(path, func(p string, i os.FileInfo, e error) error {
				if len(batch) > batchSize {
					res <- batch
					batch = nil
				}

				if e != nil {
					batch = append(batch, Resource{
						Path: p,
						Err:  e,
					})
					return e
				}

				if !isYAMLFile(i) {
					return nil
				}

				f, err := os.Open(p)
				if err != nil {
					batch = append(batch, Resource{
						Path: p,
						Err:  err,
					})
					return err
				}

				b, err := ioutil.ReadAll(f)
				if err != nil {
					batch = append(batch, Resource{
						Path: p,
						Err:  err,
					})
					return err
				}

				for _, r := range bytes.Split(b, []byte("---\n")) {
					batch = append(batch, Resource{
						Path:  p,
						Bytes: r,
					})
				}

				return nil
			})

			if len(batch) > 0 {
				res <- batch
			}
		}

		close(res)
	}()

	return res
}
