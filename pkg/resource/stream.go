package resource

import (
	"bytes"
	"io"
	"io/ioutil"
)

func FromStream(path string, r io.Reader) (<-chan Resource, <-chan error) {
	resources := make(chan Resource)
	errors := make(chan error)

	go func() {
		data, err := ioutil.ReadAll(r)
		if err != nil {
			errors <- DiscoveryError{path, err}
		}

		rawResources := bytes.Split(data, []byte("---\n"))
		for _, rawResource := range rawResources {
			resources <- Resource{Path: path, Bytes: rawResource}
		}

		close(resources)
		close(errors)
	}()

	return resources, errors
}
