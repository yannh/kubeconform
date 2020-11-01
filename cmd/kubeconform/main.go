package main

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"sync"

	"github.com/yannh/kubeconform/pkg/cache"
	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/registry"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

func downloadSchema(registries []registry.Registry, kind, version, k8sVersion string) (*gojsonschema.Schema, error) {
	var err error
	var schemaBytes []byte

	for _, reg := range registries {
		schemaBytes, err = reg.DownloadSchema(kind, version, k8sVersion)
		if err == nil {
			return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaBytes))
		}

		// If we get a 404, we try the next registry, but we exit if we get a real failure
		if er, retryable := err.(registry.Retryable); retryable && !er.IsRetryable() {
			continue
		}

		return nil, err
	}

	return nil, nil // No schema found - we don't consider it an error, resource will be skipped
}

func ValidateResources(resources <-chan resource.Resource, validationResults chan<- validator.Result, regs []registry.Registry, k8sVersion string, c *cache.SchemaCache, skip func(signature resource.Signature) bool, ignoreMissingSchemas bool) {
	for res := range resources {
		sig, err := res.Signature()
		if err != nil {
			validationResults <- validator.Result{Resource: res, Err: fmt.Errorf("error while parsing: %s", err), Status: validator.Error}
			continue
		}

		if sig.Kind == "" {
			validationResults <- validator.Result{Resource: res, Err: nil, Status: validator.Empty}
			continue // We skip resoures that don't have a Kind defined
		}

		if skip(*sig) {
			validationResults <- validator.Result{Resource: res, Err: nil, Status: validator.Skipped}
			continue
		}

		cached := false
		var schema *gojsonschema.Schema
		cacheKey := ""

		if c != nil {
			cacheKey = cache.Key(sig.Kind, sig.Version, k8sVersion)
			schema, cached = c.Get(cacheKey)
		}

		if !cached {
			if schema, err = downloadSchema(regs, sig.Kind, sig.Version, k8sVersion); err != nil {
				validationResults <- validator.Result{Resource: res, Err: err, Status: validator.Error}
				continue
			}

			if c != nil {
				c.Set(cacheKey, schema)
			}
		}

		if schema == nil {
			if ignoreMissingSchemas {
				validationResults <- validator.Result{Resource: res, Err: nil, Status: validator.Skipped}
			} else {
				validationResults <- validator.Result{Resource: res, Err: fmt.Errorf("could not find schema for %s", sig.Kind), Status: validator.Error}
			}
		}

		validationResults <- validator.Validate(res, schema)
	}
}

func processResults(o output.Output, validationResults <-chan validator.Result, result chan<- bool) {
	success := true
	for res := range validationResults {
		if res.Err != nil {
			success = false
		}
		if err := o.Write(res); err != nil {
			fmt.Fprint(os.Stderr, "failed writing log\n")
		}
	}

	result <- success
}

func realMain() int {
	var err error

	cfg := config.FromFlags()
	if cfg.Help {
		return 1
	}

	// Detect whether we have data being piped through stdin
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(cfg.Files) == 1 && cfg.Files[0] == "-" {
		isStdin = true
	}

	filter := func(signature resource.Signature) bool {
		isSkipKind, ok := cfg.SkipKinds[signature.Kind]
		return ok && isSkipKind
	}

	registries := []registry.Registry{}
	for _, schemaLocation := range cfg.SchemaLocations {
		registries = append(registries, registry.New(schemaLocation, cfg.Strict))
	}

	var o output.Output
	if o, err = output.New(cfg.OutputFormat, cfg.Summary, isStdin, cfg.Verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	var resourcesChan <-chan resource.Resource
	var errors <-chan error
	validationResults := make(chan validator.Result)
	res := make(chan bool)

	go processResults(o, validationResults, res)

	if isStdin {
		resourcesChan, errors = resource.FromStream("stdin", os.Stdin)
	} else {
		resourcesChan, errors = resource.FromFiles(cfg.Files...)
	}

	c := cache.New()
	wg := sync.WaitGroup{}
	for i := 0; i < cfg.NumberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			ValidateResources(resourcesChan, validationResults, registries, cfg.KubernetesVersion, c, filter, cfg.IgnoreMissingSchemas)
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		for err := range errors {
			if err != nil {
				if err, ok := err.(resource.DiscoveryError); ok {
					validationResults <- validator.NewError(err.Path, err.Err)
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

	close(validationResults)
	success := <-res
	o.Flush()

	if !success {
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
