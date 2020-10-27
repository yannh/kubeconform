package resource

import (
	"sigs.k8s.io/yaml"
)

type Resource struct {
	Path  string
	Bytes []byte
	Err   error
}

func (r Resource) AsMap() (map[string]interface{}, error) {
	var res map[string]interface{}
	err := yaml.Unmarshal(r.Bytes, &res)
	return res, err
}
