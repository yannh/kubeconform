package validator

import (
	"fmt"
	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"

	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

type Status int

const (
	_ Status = iota
	Error
	Skipped
	Valid
	Invalid
	Empty
)

type Validator struct {
	opts           Opts
	schemaCache    *cache.SchemaCache
	schemaDownload func(registries []registry.Registry, kind, version, k8sVersion string) (*gojsonschema.Schema, error)
	regs           []registry.Registry
}

type Opts struct {
	SkipTLS              bool
	SkipKinds            map[string]bool
	RejectKinds          map[string]bool
	KubernetesVersion    string
	Strict               bool
	IgnoreMissingSchemas bool
}

func downloadSchema(registries []registry.Registry, kind, version, k8sVersion string) (*gojsonschema.Schema, error) {
	var err error
	var schemaBytes []byte

	for _, reg := range registries {
		schemaBytes, err = reg.DownloadSchema(kind, version, k8sVersion)
		if err == nil {
			return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaBytes))
		}

		// If we get a 404, we try the next registry, but we exit if we get a real failure
		if _, notfound := err.(*registry.NotFoundError); notfound {
			continue
		}

		return nil, err
	}

	return nil, nil // No schema found - we don't consider it an error, resource will be skipped
}

func New(schemaLocations []string, opts Opts) *Validator {
	registries := []registry.Registry{}
	for _, schemaLocation := range schemaLocations {
		registries = append(registries, registry.New(schemaLocation, opts.Strict, opts.SkipTLS))
	}

	if opts.SkipKinds == nil {
		opts.SkipKinds = map[string]bool{}
	}
	if opts.RejectKinds == nil {
		opts.RejectKinds = map[string]bool{}
	}

	return &Validator{
		opts:           opts,
		schemaDownload: downloadSchema,
		schemaCache:    cache.New(),
		regs:           registries,
	}
}

func (v *Validator) Validate(res resource.Resource) Result {
	skip := func(signature resource.Signature) bool {
		isSkipKind, ok := v.opts.SkipKinds[signature.Kind]
		return ok && isSkipKind
	}

	reject := func(signature resource.Signature) bool {
		_, ok := v.opts.RejectKinds[signature.Kind]
		return ok
	}

	sig, err := res.Signature()
	if err != nil {
		return Result{Resource: res, Err: fmt.Errorf("error while parsing: %s", err), Status: Error}
	}

	if sig.Kind == "" {
		return Result{Resource: res, Err: nil, Status: Empty}
	}

	if skip(*sig) {
		return Result{Resource: res, Err: nil, Status: Skipped}
	}

	if reject(*sig) {
		return Result{Resource: res, Err: fmt.Errorf("prohibited resource kind %s", sig.Kind), Status: Error}
	}

	cached := false
	var schema *gojsonschema.Schema
	cacheKey := ""

	if v.schemaCache != nil {
		cacheKey = cache.Key(sig.Kind, sig.Version, v.opts.KubernetesVersion)
		schema, cached = v.schemaCache.Get(cacheKey)
	}

	if !cached {
		if schema, err = v.schemaDownload(v.regs, sig.Kind, sig.Version, v.opts.KubernetesVersion); err != nil {
			return Result{Resource: res, Err: err, Status: Error}
		}

		if v.schemaCache != nil {
			v.schemaCache.Set(cacheKey, schema)
		}
	}

	if schema == nil {
		if v.opts.IgnoreMissingSchemas {
			return Result{Resource: res, Err: nil, Status: Skipped}
		} else {
			return Result{Resource: res, Err: fmt.Errorf("could not find schema for %s", sig.Kind), Status: Error}
		}
	}

	if schema == nil {
		return Result{Resource: res, Status: Skipped, Err: nil}
	}

	var resource map[string]interface{}
	if err := yaml.Unmarshal(res.Bytes, &resource); err != nil {
		return Result{Resource: res, Status: Error, Err: fmt.Errorf("error unmarshalling resource: %s", err)}
	}
	resourceLoader := gojsonschema.NewGoLoader(resource)

	results, err := schema.Validate(resourceLoader)
	if err != nil {
		// This error can only happen if the Object to validate is poorly formed. There's no hope of saving this one
		return Result{Resource: res, Status: Error, Err: fmt.Errorf("problem validating schema. Check JSON formatting: %s", err)}
	}

	if results.Valid() {
		return Result{Resource: res, Status: Valid}
	}

	msg := ""
	for _, errMsg := range results.Errors() {
		if msg != "" {
			msg += " - "
		}
		msg += errMsg.Description()
	}

	return Result{Resource: res, Status: Invalid, Err: fmt.Errorf("%s", msg)}
}

// ValidFormat is a type for quickly forcing
// new formats on the gojsonschema loader
type ValidFormat struct{}

// ValidFormat is a type for quickly forcing
// new formats on the gojsonschema loader
func (f ValidFormat) IsFormat(input interface{}) bool {
	return true
}

// From kubeval - let's see if absolutely necessary
// func init () {
// 	gojsonschema.FormatCheckers.Add("int64", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("byte", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("int32", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("int-or-string", ValidFormat{})
// }

// Result contains the details of the result of a resource validation
type Result struct {
	Resource resource.Resource
	Err      error
	Status   Status
}

// NewError is a utility function to generate a validation error
func NewError(filename string, err error) Result {
	return Result{
		Resource: resource.Resource{Path: filename},
		Err:      err,
		Status:   Error,
	}
}
