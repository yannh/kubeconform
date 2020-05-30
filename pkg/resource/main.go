package resource

import (
	"sigs.k8s.io/yaml"
)

type Signature struct {
	Kind, Version, Namespace string
}

// TODO: Support multi-resources yaml files
func SignatureFromBytes(s []byte) (Signature, error) {
	resource := struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Namespace string `yaml:"Namespace"`
		} `yaml:"Metadata"`
	}{}
	err := yaml.Unmarshal(s, &resource)

	return Signature{Kind: resource.Kind, Version: resource.APIVersion, Namespace: resource.Metadata.Namespace}, err
}

