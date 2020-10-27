package schema

import (
	"errors"
	"fmt"
	"sync"
)

// Fetcher TODO
type Fetcher interface {
	Get(kind string, version string, kubernetesVersion string) (*Schema, error)
}

// Repository TODO
type Repository struct {
	schemas     map[string]*Schema
	schemasLock sync.RWMutex

	fetcher []Fetcher
}

func key(resourceKind, resourceAPIVersion, k8sVersion string) string {
	return fmt.Sprintf("%s-%s-%s", resourceKind, resourceAPIVersion, k8sVersion)
}

// Get TODO
func (r *Repository) Get(kind string, version string, kubernetesVersion string) (*Schema, error) {
	r.schemasLock.RLock()
	defer r.schemasLock.RUnlock()

	schema, ok := r.schemas[key(kind, version, kubernetesVersion)]
	if ok {
		return schema, nil
	}

	for _, fetcher := range r.fetcher {
		schema, err := fetcher.Get(kind, version, kubernetesVersion)

		if err != nil {
			continue
		}

		r.schemas[key(kind, version, kubernetesVersion)] = schema
		return schema, nil
	}

	return nil, errors.New("schema not found")
}

// Option TODO
type Option func(*Repository)

// New TODO
func New(opts ...Option) *Repository {
	r := &Repository{
		schemas:     map[string]*Schema{},
		schemasLock: sync.RWMutex{},
		fetcher:     []Fetcher{},
	}

	for _, opt := range opts {
		opt(r)
	}

	// add kubernetesjsonschema.dev as last fetcher
	FromRemote(kubernetesJSONSchemaURLTmpl)(r)

	return r
}
