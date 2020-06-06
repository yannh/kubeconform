package resource

import (
	"sigs.k8s.io/yaml"
)

type Signature struct {
	Kind, Version, Namespace string
}

// SignatureFromBytes returns key identifying elements from a []byte representing the resource
func SignatureFromBytes(res []byte) (Signature, error) {
	resource := struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Namespace string `yaml:"Namespace"`
		} `yaml:"Metadata"`
	}{}
	err := yaml.Unmarshal(res, &resource)

	return Signature{Kind: resource.Kind, Version: resource.APIVersion, Namespace: resource.Metadata.Namespace}, err
}
