package resource

import (
	"sigs.k8s.io/yaml"
)

type Signature struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name         string `yaml:"name"`
		Namespace    string `yaml:"namespace"`
		GenerateName string `yaml:"generateName"`
	} `yaml:"Metadata"`
}

func (r Resource) Signature() (Signature, error) {
	var sig Signature
	err := yaml.Unmarshal(r.Bytes, &sig)
	return sig, err
}
