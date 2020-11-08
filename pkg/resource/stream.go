package resource

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
)

func FromStream(ctx context.Context, path string, r io.Reader) (<-chan Resource, <-chan error) {
	resources := make(chan Resource)
	errors := make(chan error)
	stop := false

	go func() {
		<-ctx.Done()
		stop = true
	}()

	go func() {
		data, err := ioutil.ReadAll(r)
		if err != nil {
			errors <- DiscoveryError{path, err}
		}

		rawResources := bytes.Split(data, []byte("---\n"))
		for _, rawResource := range rawResources {
			if stop == true {
				break
			}
			resources <- Resource{Path: path, Bytes: rawResource}
		}

		close(resources)
		close(errors)
	}()

	return resources, errors
}
