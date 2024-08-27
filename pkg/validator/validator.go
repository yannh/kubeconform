// This is the main package to import to embed kubeconform in your software
package validator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"sigs.k8s.io/yaml"
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

	return &v{
		opts:           opts,
		schemaDownload: downloadSchema,
		schemaCache:    cache.NewInMemoryCache(),
		regs:           registries,
	}, nil
}

type v struct {
	opts           Opts
	schemaCache    cache.Cache
	schemaDownload func(registries []registry.Registry, kind, version, k8sVersion string) (*jsonschema.Schema, error)
	regs           []registry.Registry
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

	if val.schemaCache != nil {
		s, err := val.schemaCache.Get(sig.Kind, sig.Version, val.opts.KubernetesVersion)
		if err == nil {
			cached = true
			schema = s.(*jsonschema.Schema)
		}
	}

	if !cached {
		if schema, err = val.schemaDownload(val.regs, sig.Kind, sig.Version, val.opts.KubernetesVersion); err != nil {
			return Result{Resource: res, Err: err, Status: Error}
		}

		if val.schemaCache != nil {
			val.schemaCache.Set(sig.Kind, sig.Version, val.opts.KubernetesVersion, schema)
		}
	}

	if schema == nil {
		if val.opts.IgnoreMissingSchemas {
			return Result{Resource: res, Err: nil, Status: Skipped}
		}

		return Result{Resource: res, Err: fmt.Errorf("could not find schema for %s", sig.Kind), Status: Error}
	}

	validationErrors := []ValidationError{}
	err = schema.Validate(r)
	if err != nil {
		var e *jsonschema.ValidationError
		if errors.As(err, &e) {
			for _, ve := range e.Causes {
				validationErrors = append(validationErrors, ValidationError{
					Path: ve.InstanceLocation,
					Msg:  ve.Message,
				})
			}
		}
		return Result{
			Resource:         res,
			Status:           Invalid,
			Err:              fmt.Errorf("problem validating schema. Check JSON formatting: %s", err),
			ValidationErrors: validationErrors,
		}
	}

	if res.Metadata != nil {
	    metadataPath := res.Path + " - .metadata"
	    namePath := metadataPath + ".name"
	    name := res.Metadata.Name
	    if name == "" {
	        if res.Metadata.GenerateName != "" {
                name =  res.Metadata.GenerateName
	            namePath = metadataPath + ".generateName"
            }
	    }

	    if !validateDnsLabels(name) {
            validationErrors = append(validationErrors, ValidationError{
                    Path: namePath,
                    Msg:  "invalid metadata name",
                })
        }

		validationErrors = validateMeta(metadataPath, res.Metadata, validationErrors)
        if len(validationErrors) > 0 {
            return Result{
                        Resource:         res,
                        Status:           Invalid,
                        Err:              fmt.Errorf("invalid metadata."),
                        ValidationErrors: validationErrors,
                    }
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

func downloadSchema(registries []registry.Registry, kind, version, k8sVersion string) (*jsonschema.Schema, error) {
	var err error
	var schemaBytes []byte
	var path string

	for _, reg := range registries {
		path, schemaBytes, err = reg.DownloadSchema(kind, version, k8sVersion)
		if err == nil {
			c := jsonschema.NewCompiler()
			c.Draft = jsonschema.Draft4
			if err := c.AddResource(path, bytes.NewReader(schemaBytes)); err != nil {
				continue
			}
			schema, err := c.Compile(path)
			// If we got a non-parseable response, we try the next registry
			if err != nil {
				continue
			}
			return schema, err
		}

		// If we get a 404, we try the next registry, but we exit if we get a real failure
		if _, notfound := err.(*registry.NotFoundError); notfound {
			continue
		}

		return nil, err
	}

	return nil, nil // No schema found - we don't consider it an error, resource will be skipped
}

func validateMeta(path string, metadata *resource.ObjectMeta, validationErrors []ValidationError) []ValidationError {
	if metadata.Annotations != nil {
		validationErrors = validateAnnotations(path + ".annotations", metadata.Annotations, validationErrors)
	}
	if metadata.Labels != nil {
		validationErrors = validateLabels(path + ".labels", metadata.Labels, validationErrors)
	}
    return validationErrors
}

/* Annotations are key/value pairs. */
func validateAnnotations(path string, annotations map[string]string, validationErrors []ValidationError) []ValidationError {
    for k, v := range annotations {
        keypath := path + "[" + k + "]"
        validationErrors = validateKey(keypath, k, validationErrors)
        if !validateAnnotationValue(v) {
            validationErrors = append(validationErrors, ValidationError{
                                Path: keypath,
                                Msg:  "invalid annotation value",
                            })
        }
    }
    return validationErrors
}

/* Labels are key/value pairs.
 */
func validateLabels(path string, labels map[string]string, validationErrors []ValidationError) []ValidationError {
    for k, v := range labels {
        keypath := path + "[" + k + "]"
        validationErrors= validateKey(keypath, k, validationErrors)
        if !validateNameSegment(v) {
            validationErrors = append(validationErrors, ValidationError{
                                            Path: keypath,
                                            Msg:  "invalid label value",
                                        })
        }
    }
    return validationErrors
}

/* Valid keys have two segments: an optional prefix and name, separated by a slash (/)
*/
func validateKey(keypath string, key string, validationErrors []ValidationError) []ValidationError {
    if len(key) == 0 {
        validationErrors = append(validationErrors, ValidationError{
                                        Path: keypath,
                                        Msg:  "invalid annotation key",
                                    })
    } else {
        var name string
        prefix, suffix, found := strings.Cut(key, "/")
        if found {
            name = suffix
            if !validateDnsLabels(prefix) {
                validationErrors = append(validationErrors, ValidationError{
                                                Path: keypath,
                                                Msg:  "invalid annotation key prefix",
                                            })
            }
        } else {
            name = key
        }

        if !validateNameSegment(name) {
            validationErrors = append(validationErrors, ValidationError{
                                            Path: keypath,
                                            Msg:  "invalid annotation key name",
                                        })

        }
    }
    return validationErrors
}

var alphanumericPlusUnderscorePeriodHyphen = regexp.MustCompile("^[0-9A-Za-z_.-]+$")

func isAlphaNumeric(v byte) bool {
  return (v >= '0' && v <= '9') ||
         (v >= 'A' && v <= 'Z') ||
         (v >= 'a' && v <= 'z')
}

/* The name segment must be 63 characters or less, beginning and ending with an alphanumeric character
   ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.
*/
func validateNameSegment(name string) bool {
    return len(name) <= 63 &&
        alphanumericPlusUnderscorePeriodHyphen.MatchString(name) &&
        isAlphaNumeric(name[0]) &&
        isAlphaNumeric(name[len(name)-1])
}

var alphanumericPlusHyphen = regexp.MustCompile("^[0-9A-Za-z-]+$")

/* The domain name may not exceed the length of 253 characters in its textual representation.
   A label may contain one to 63 characters of a through z, A through Z, digits 0 through 9, and hyphen.
   Labels may not start or end with a hyphen.
*/
func validateDnsLabels(domain string) bool {
    if len(domain) == 0 || len(domain) > 253 {
        return false
    } else {
        labels := strings.Split(domain, ".")
        for _, label := range labels {
            if len(label) == 0 ||
                    len(label) > 63 ||
                    !alphanumericPlusHyphen.MatchString(label) ||
                    label[0] == '-' ||
                    label[len(label)-1] == '-' {
                return false
            }
        }
    }
    return true
}

/* annotation must have value
*/
func validateAnnotationValue(value string) bool {
    return len(value) != 0
}
