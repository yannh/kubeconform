package resource

import (
	"sigs.k8s.io/yaml"
)

type Signature struct {
	Kind, Version, Namespace, Name string
}

// SignatureFromBytes returns key identifying elements from a []byte representing the resource
func SignatureFromBytes(res []byte) (Signature, error) {
	resource := struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Name      string `yaml:"Name"`
			Namespace string `yaml:"Namespace"`
		} `yaml:"Metadata"`
	}{}
	err := yaml.Unmarshal(res, &resource)

	return Signature{Kind: resource.Kind, Version: resource.APIVersion, Namespace: resource.Metadata.Namespace, Name: resource.Metadata.Name}, err
}
