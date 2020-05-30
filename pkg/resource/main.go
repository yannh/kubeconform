package resource

import (
	"io"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
)

type Resource struct {
	Kind, Version, Namespace string
}

// TODO: Support multi-resources yaml files
func Read(r io.Reader) (Resource, error) {
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return Resource{}, err
	}

	resource := struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Namespace string `yaml:"Namespace"`
		} `yaml:"Metadata"`
	}{}
	err = yaml.Unmarshal(s, &resource)

	return Resource{Kind: resource.Kind, Version: resource.APIVersion, Namespace: resource.Metadata.Namespace}, err
}