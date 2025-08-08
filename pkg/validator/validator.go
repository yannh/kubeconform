// This is the main package to import to embed kubeconform in your software
package validator

import (
	"context"
	"errors"
	"fmt"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/loader"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

// Different types of validation results
type Status int

const (
	_       Status = iota
	Error          // an error occurred processing the file / resource
	Skipped        // resource has been skipped, for example if its Kind was part of the kinds to skip
	Valid          // resource is valid
	Invalid        // resource is invalid
	Empty          // resource is empty. Note: is triggered for files starting with a --- separator.
)

type ValidationError struct {
	Path string `json:"path"`
	Msg  string `json:"msg"`
}

func (ve *ValidationError) Error() string {
	return ve.Msg
}

// Result contains the details of the result of a resource validation
type Result struct {
	Resource         resource.Resource
	Err              error
	Status           Status
	ValidationErrors []ValidationError
}

// Validator exposes multiple methods to validate your Kubernetes resources.
type Validator interface {
	ValidateResource(res resource.Resource) Result
	Validate(filename string, r io.ReadCloser) []Result
	ValidateWithContext(ctx context.Context, filename string, r io.ReadCloser) []Result
}

// Opts contains a set of options for the validator.
type Opts struct {
	Cache                string              // Cache schemas downloaded via HTTP to this folder
	Debug                bool                // Debug infos will be print here
	SkipTLS              bool                // skip TLS validation when downloading from an HTTP Schema Registry
	SkipKinds            map[string]struct{} // List of resource Kinds to ignore
	RejectKinds          map[string]struct{} // List of resource Kinds to reject
	KubernetesVersion    string              // Kubernetes Version - has to match one in https://github.com/instrumenta/kubernetes-json-schema
	Strict               bool                // thros an error if resources contain undocumented fields
	IgnoreMissingSchemas bool                // skip a resource if no schema for that resource can be found
}

// New returns a new Validator
func New(schemaLocations []string, opts Opts) (Validator, error) {
	// Default to our kubernetes-json-schema fork
	// raw.githubusercontent.com is frontend by Fastly and very fast
	if len(schemaLocations) == 0 {
		schemaLocations = []string{"https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json"}
	}

	registries := []registry.Registry{}
	for _, schemaLocation := range schemaLocations {
		reg, err := registry.New(schemaLocation, opts.Cache, opts.Strict, opts.SkipTLS, opts.Debug)
		if err != nil {
			return nil, err
		}
		registries = append(registries, reg)
	}

	if opts.KubernetesVersion == "" {
		opts.KubernetesVersion = "master"
	}

	if opts.SkipKinds == nil {
		opts.SkipKinds = map[string]struct{}{}
	}
	if opts.RejectKinds == nil {
		opts.RejectKinds = map[string]struct{}{}
	}

	var filecache cache.Cache = nil
	if opts.Cache != "" {
		fi, err := os.Stat(opts.Cache)
		if err != nil {
			return nil, fmt.Errorf("failed opening cache folder %s: %s", opts.Cache, err)
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("cache folder %s is not a directory", err)
		}

		filecache = cache.NewOnDiskCache(opts.Cache)
	}

	httpLoader, err := loader.NewHTTPURLLoader(false, filecache)
	if err != nil {
		return nil, fmt.Errorf("failed creating HTTP loader: %s", err)
	}

	return &v{
		opts:              opts,
		schemaDownload:    downloadSchema,
		schemaMemoryCache: cache.NewInMemoryCache(),
		regs:              registries,
		loader: jsonschema.SchemeURLLoader{
			"file":  jsonschema.FileLoader{},
			"http":  httpLoader,
			"https": httpLoader,
		},
	}, nil
}

type v struct {
	opts              Opts
	schemaDiskCache   cache.Cache
	schemaMemoryCache cache.Cache
	schemaDownload    func(registries []registry.Registry, loader jsonschema.SchemeURLLoader, kind, version, k8sVersion string) (*jsonschema.Schema, error)
	regs              []registry.Registry
	loader            jsonschema.SchemeURLLoader
}

func key(resourceKind, resourceAPIVersion, k8sVersion string) string {
	return fmt.Sprintf("%s-%s-%s", resourceKind, resourceAPIVersion, k8sVersion)
}

// ValidateResource validates a single resource. This allows to validate
// large resource streams using multiple Go Routines.
func (val *v) ValidateResource(res resource.Resource) Result {
	// For backward compatibility reasons when determining whether
	// a resource should be skipped or rejected we use both
	// the GVK encoding of the resource signatures (the recommended method
	// for skipping/rejecting resources) and the raw Kind.

	skip := func(signature resource.Signature) bool {
		if _, ok := val.opts.SkipKinds[signature.GroupVersionKind()]; ok {
			return ok
		}
		_, ok := val.opts.SkipKinds[signature.Kind]
		return ok
	}

	reject := func(signature resource.Signature) bool {
		if _, ok := val.opts.RejectKinds[signature.GroupVersionKind()]; ok {
			return ok
		}
		_, ok := val.opts.RejectKinds[signature.Kind]
		return ok
	}

	if len(res.Bytes) == 0 {
		return Result{Resource: res, Err: nil, Status: Empty}
	}

	var r map[string]interface{}
	unmarshaller := yaml.Unmarshal
	if val.opts.Strict {
		unmarshaller = yaml.UnmarshalStrict
	}

	if err := unmarshaller(res.Bytes, &r); err != nil {
		return Result{Resource: res, Status: Error, Err: fmt.Errorf("error unmarshalling resource: %s", err)}
	}

	if r == nil { // Resource is empty
		return Result{Resource: res, Err: nil, Status: Empty}
	}

	sig, err := res.SignatureFromMap(r)
	if err != nil {
		return Result{Resource: res, Err: fmt.Errorf("error while parsing: %s", err), Status: Error}
	}

	if skip(*sig) {
		return Result{Resource: res, Err: nil, Status: Skipped}
	}

	if reject(*sig) {
		return Result{Resource: res, Err: fmt.Errorf("prohibited resource kind %s", sig.Kind), Status: Error}
	}

	cached := false
	var schema *jsonschema.Schema

	if val.schemaMemoryCache != nil {
		s, err := val.schemaMemoryCache.Get(key(sig.Kind, sig.Version, val.opts.KubernetesVersion))
		if err == nil {
			cached = true
			schema = s.(*jsonschema.Schema)
		}
	}

	if !cached {
		if schema, err = val.schemaDownload(val.regs, val.loader, sig.Kind, sig.Version, val.opts.KubernetesVersion); err != nil {
			return Result{Resource: res, Err: err, Status: Error}
		}

		if val.schemaMemoryCache != nil {
			val.schemaMemoryCache.Set(key(sig.Kind, sig.Version, val.opts.KubernetesVersion), schema)
		}
	}

	if schema == nil {
		if val.opts.IgnoreMissingSchemas {
			return Result{Resource: res, Err: nil, Status: Skipped}
		}

		return Result{Resource: res, Err: fmt.Errorf("could not find schema for %s", sig.Kind), Status: Error}
	}

	err = schema.Validate(r)
	if err != nil {
		validationErrors := []ValidationError{}
		var e *jsonschema.ValidationError
		if errors.As(err, &e) {
			for _, ve := range e.Causes {
				path := ""
				for _, f := range ve.InstanceLocation {
					path = path + "/" + f
				}
				validationErrors = append(validationErrors, ValidationError{
					Path: path,
					Msg:  ve.ErrorKind.LocalizedString(message.NewPrinter(language.English)),
				})
			}

		}

		return Result{
			Resource:         res,
			Status:           Invalid,
			Err:              fmt.Errorf("problem validating schema. Check JSON formatting: %s", strings.ReplaceAll(err.Error(), "\n", " ")),
			ValidationErrors: validationErrors,
		}
	}

	return Result{Resource: res, Status: Valid}
}

// ValidateWithContext validates resources found in r
// filename should be a name for the stream, such as a filename or stdin
func (val *v) ValidateWithContext(ctx context.Context, filename string, r io.ReadCloser) []Result {
	validationResults := []Result{}
	resourcesChan, _ := resource.FromStream(ctx, filename, r)
	for {
		select {
		case res, ok := <-resourcesChan:
			if ok {
				validationResults = append(validationResults, val.ValidateResource(res))
			} else {
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

// validateDuration is a custom validator for the duration format
// as JSONSchema only supports the ISO 8601 format, i.e. `PT1H30M`,
// while Kubernetes API machinery expects the Go duration format, i.e. `1h30m`
// which is commonly used in Kubernetes operators for specifying intervals.
// https://github.com/kubernetes/apiextensions-apiserver/blob/1ecd29f74da0639e2e6e3b8fac0c9bfd217e05eb/pkg/apis/apiextensions/v1/types_jsonschema.go#L71
func validateDuration(v any) error {
	// Try validation with the Go duration format
	if _, err := time.ParseDuration(v.(string)); err == nil {
		return nil
	}

	s, ok := v.(string)
	if !ok {
		return nil
	}

	// must start with 'P'
	s, ok = strings.CutPrefix(s, "P")
	if !ok {
		return fmt.Errorf("must start with P")
	}
	if s == "" {
		return fmt.Errorf("nothing after P")
	}

	// dur-week
	if s, ok := strings.CutSuffix(s, "W"); ok {
		if s == "" {
			return fmt.Errorf("no number in week")
		}
		for _, ch := range s {
			if ch < '0' || ch > '9' {
				return fmt.Errorf("invalid week")
			}
		}
		return nil
	}

	allUnits := []string{"YMD", "HMS"}
	for i, s := range strings.Split(s, "T") {
		if i != 0 && s == "" {
			return fmt.Errorf("no time elements")
		}
		if i >= len(allUnits) {
			return fmt.Errorf("more than one T")
		}
		units := allUnits[i]
		for s != "" {
			digitCount := 0
			for _, ch := range s {
				if ch >= '0' && ch <= '9' {
					digitCount++
				} else {
					break
				}
			}
			if digitCount == 0 {
				return fmt.Errorf("missing number")
			}
			s = s[digitCount:]
			if s == "" {
				return fmt.Errorf("missing unit")
			}
			unit := s[0]
			j := strings.IndexByte(units, unit)
			if j == -1 {
				if strings.IndexByte(allUnits[i], unit) != -1 {
					return fmt.Errorf("unit %q out of order", unit)
				}
				return fmt.Errorf("invalid unit %q", unit)
			}
			units = units[j+1:]
			s = s[1:]
		}
	}

	return nil
}

func downloadSchema(registries []registry.Registry, l jsonschema.SchemeURLLoader, kind, version, k8sVersion string) (*jsonschema.Schema, error) {
	var err error
	var path string
	var s any

	for _, reg := range registries {
		path, s, err = reg.DownloadSchema(kind, version, k8sVersion)
		if err == nil {
			c := jsonschema.NewCompiler()
			c.RegisterFormat(&jsonschema.Format{"duration", validateDuration})
			c.UseLoader(l)
			c.DefaultDraft(jsonschema.Draft4)
			if err := c.AddResource(path, s); err != nil {
				continue
			}
			schema, err := c.Compile(path)
			// If we got a non-parseable response, we try the next registry
			if err != nil {
				continue
			}
			if schema != nil {
				// First check if the schema has explicit kind constraints
				if hasKindConstraints(schema) {
					if isSchemaForKind(schema, kind) {
						return schema, nil
					}
				} else {
					// If no kind constraints, check if the path matches the kind
					if pathMatchesKind(path, kind) {
						return schema, nil
					}
				}
			}
			continue
		}

		if _, notfound := err.(*loader.NotFoundError); notfound {
			continue
		}
		if _, nonJSONError := err.(*loader.NonJSONResponseError); nonJSONError {
			continue
		}

		return nil, err
	}

	return nil, nil // No schema found - we don't consider it an error, resource will be skipped
}

func hasKindConstraints(schema *jsonschema.Schema) bool {
	if schema == nil {
		return false
	}

	if schema.Properties != nil {
		if kindProp, ok := schema.Properties["kind"]; ok {
			// Check if the kind property has enum or const constraints
			if kindProp.Enum != nil || kindProp.Const != nil {
				return true
			}
		}
	}

	return false
}

func pathMatchesKind(path string, expectedKind string) bool {
	// Convert expected kind to lowercase for comparison
	expectedLower := strings.ToLower(expectedKind)

	// Check if the path contains the expected kind name
	pathLower := strings.ToLower(path)

	// For direct file paths (like secretstore_v1.json), check if the filename contains the kind
	if strings.Contains(pathLower, expectedLower) {
		return true
	}

	return false
}

func isSchemaForKind(schema *jsonschema.Schema, expectedKind string) bool {
	if schema == nil {
		return true
	}

	if schema.Properties != nil {
		if kindProp, ok := schema.Properties["kind"]; ok {
			if kindProp.Enum != nil {
				if enumValues := kindProp.Enum.Values; len(enumValues) > 0 {
					if enumKind, ok := enumValues[0].(string); ok {
						return enumKind == expectedKind
					}
				}
			}
			if kindProp.Const != nil {
				if constVal, ok := (*kindProp.Const).(string); ok {
					return constVal == expectedKind
				}
			}
		}
	}

	return true
}
