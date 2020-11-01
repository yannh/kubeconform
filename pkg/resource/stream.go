package resource

import (
	"bytes"
	"io"
	"io/ioutil"
)

func FromStream(path string, r io.Reader) ([]Resource, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return []Resource{}, err
	}

	resources := []Resource{}
	rawResources := bytes.Split(data, []byte("---\n"))
	for _, rawResource := range rawResources {
		resources = append(resources, Resource{Path: path, Bytes: rawResource})
	}

	return resources, nil
}
