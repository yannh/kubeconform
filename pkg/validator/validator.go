package validator

import (
	"context"
	"fmt"
	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"io"

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

// Result contains the details of the result of a resource validation
type Result struct {
	Resource resource.Resource
	Err      error
	Status   Status
}

type Validator interface {
	ValidateResource(res resource.Resource) Result
	Validate(filename string, r io.ReadCloser) []Result
	ValidateWithContext(ctx context.Context, filename string, r io.ReadCloser) []Result
}

// Opts contains a set of options for the validator.
type Opts struct {
	SkipTLS              bool            // skip TLS validation when downloading from an HTTP Schema Registry
	SkipKinds            map[string]bool // List of resource Kinds to ignore
	RejectKinds          map[string]bool // List of resource Kinds to reject
	KubernetesVersion    string          // Kubernetes Version - has to match one in https://github.com/instrumenta/kubernetes-json-schema
	Strict               bool            // thros an error if resources contain undocumented fields
	IgnoreMissingSchemas bool            // skip a resource if no schema for that resource can be found
}

// New returns a new Validator
func New(schemaLocations []string, opts Opts) Validator {
	// Default to kubernetesjsonschema.dev
	if schemaLocations == nil || len(schemaLocations) == 0 {
		schemaLocations = []string{"https://kubernetesjsonschema.dev"}
	}

	registries := []registry.Registry{}
	for _, schemaLocation := range schemaLocations {
		registries = append(registries, registry.New(schemaLocation, opts.Strict, opts.SkipTLS))
	}

	if opts.KubernetesVersion == "" {
		opts.KubernetesVersion = "1.18.0"
	}

	if opts.SkipKinds == nil {
		opts.SkipKinds = map[string]bool{}
	}
	if opts.RejectKinds == nil {
		opts.RejectKinds = map[string]bool{}
	}

	return &v{
		opts:           opts,
		schemaDownload: downloadSchema,
		schemaCache:    cache.New(),
		regs:           registries,
	}
}

type v struct {
	opts           Opts
	schemaCache    *cache.SchemaCache
	schemaDownload func(registries []registry.Registry, kind, version, k8sVersion string) (*gojsonschema.Schema, error)
	regs           []registry.Registry
}

// ValidateResource validates a single resource. This allows to validate
// large resource streams using multiple Go Routines.
func (val *v) ValidateResource(res resource.Resource) Result {
	skip := func(signature resource.Signature) bool {
		isSkipKind, ok := val.opts.SkipKinds[signature.Kind]
		return ok && isSkipKind
	}

	reject := func(signature resource.Signature) bool {
		_, ok := val.opts.RejectKinds[signature.Kind]
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

	if val.schemaCache != nil {
		cacheKey = cache.Key(sig.Kind, sig.Version, val.opts.KubernetesVersion)
		schema, cached = val.schemaCache.Get(cacheKey)
	}

	if !cached {
		if schema, err = val.schemaDownload(val.regs, sig.Kind, sig.Version, val.opts.KubernetesVersion); err != nil {
			return Result{Resource: res, Err: err, Status: Error}
		}

		if val.schemaCache != nil {
			val.schemaCache.Set(cacheKey, schema)
		}
	}

	if schema == nil {
		if val.opts.IgnoreMissingSchemas {
			return Result{Resource: res, Err: nil, Status: Skipped}
		} else {
			return Result{Resource: res, Err: fmt.Errorf("could not find schema for %s", sig.Kind), Status: Error}
		}
	}

	var r map[string]interface{}
	if err := yaml.Unmarshal(res.Bytes, &r); err != nil {
		return Result{Resource: res, Status: Error, Err: fmt.Errorf("error unmarshalling resource: %s", err)}
	}
	resourceLoader := gojsonschema.NewGoLoader(r)

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

// ValidateWithContext validates resources found in r
// filename should be a name for the stream, such as a filename or stdin
func (val *v) ValidateWithContext(ctx context.Context, filename string, r io.ReadCloser) []Result {
	validationResults := []Result{}
	resourcesChan, _ := resource.FromStream(ctx, filename, r)
	for {
		select {
		case res, ok := <-resourcesChan:
			validationResults = append(validationResults, val.ValidateResource(res))
			if !ok {
				resourcesChan = nil
			}

		case <-ctx.Done():
			break
		}

		if resourcesChan == nil {
			break
		}
	}

	r.Close()
	return validationResults
}

// Validate validates resources found in r
// filename should be a name for the stream, such as a filename or stdin
func (val *v) Validate(filename string, r io.ReadCloser) []Result {
	return val.ValidateWithContext(context.Background(), filename, r)
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

// From kubeval - let's see if absolutely necessary
// func init () {
// 	gojsonschema.FormatCheckers.Add("int64", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("byte", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("int32", ValidFormat{})
// 	gojsonschema.FormatCheckers.Add("int-or-string", ValidFormat{})
// }
