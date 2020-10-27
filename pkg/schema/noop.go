package schema

import "errors"

type noop struct{}

func (n noop) Get(kind string, version, kubernetesVersion string) (*Schema, error) {
	return nil, errors.New("schema not found")
}
