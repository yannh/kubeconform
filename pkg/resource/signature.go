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
			Name         string `yaml:"name"`
			Namespace    string `yaml:"namespace"`
			GenerateName string `yaml:"generateName"`
		} `yaml:"Metadata"`
	}{}
	err := yaml.Unmarshal(res, &resource)

	name := resource.Metadata.Name
	if resource.Metadata.GenerateName != "" {
		name = resource.Metadata.GenerateName + "{{ generateName }}"
	}

	return Signature{Kind: resource.Kind, Version: resource.APIVersion, Namespace: resource.Metadata.Namespace, Name: name}, err
}
