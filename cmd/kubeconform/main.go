package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"sigs.k8s.io/yaml"

	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

var version = "development"

// New function to load the manifest and schema for injecting defaults
func loadManifestAndSchema(manifestPath, schemaPath string) (map[string]interface{}, map[string]interface{}, error) {
	manifestFile, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading manifest file: %v", err)
	}
	var manifest map[string]interface{}
	if err := yaml.Unmarshal(manifestFile, &manifest); err != nil {
		return nil, nil, fmt.Errorf("error parsing manifest YAML: %v", err)
	}

	schemaFile, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading schema file: %v", err)
	}
	var schema map[string]interface{}
	if err := json.Unmarshal(schemaFile, &schema); err != nil {
		return nil, nil, fmt.Errorf("error parsing schema JSON: %v", err)
	}

	return manifest, schema, nil
}

// New function to inject defaults recursively
func injectDefaultsRecursively(schema map[string]interface{}, manifest map[string]interface{}) {
	properties, propertiesExist := schema["properties"].(map[string]interface{})
	if !propertiesExist {
		return
	}

	for key, subschema := range properties {
		if _, keyExists := manifest[key]; !keyExists {
			subSchemaMap, ok := subschema.(map[string]interface{})
			if ok {
				if defaultValue, hasDefault := subSchemaMap["default"]; hasDefault {
					manifest[key] = defaultValue
					fmt.Printf("Injected default for %s: %v\\n", key, defaultValue)
				}
			}
		} else {
			if subSchemaMap, ok := subschema.(map[string]interface{}); ok {
				if subSchemaType, typeExists := subSchemaMap["type"].(string); typeExists {
					if subSchemaType == "object" {
						if nestedManifest, isMap := manifest[key].(map[string]interface{}); isMap {
							injectDefaultsRecursively(subSchemaMap, nestedManifest)
						}
					} else if subSchemaType == "array" {
						if arrayItems, hasItems := subSchemaMap["items"].(map[string]interface{}); hasItems {
							if manifestArray, isArray := manifest[key].([]interface{}); isArray {
								for _, item := range manifestArray {
									if itemMap, isMap := item.(map[string]interface{}); isMap {
										injectDefaultsRecursively(arrayItems, itemMap)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func processResults(cancel context.CancelFunc, o output.Output, validationResults <-chan validator.Result, exitOnError bool) <-chan bool {
	success := true
	result := make(chan bool)

	go func() {
		for res := range validationResults {
			if res.Status == validator.Error || res.Status == validator.Invalid {
				success = false
			}
			if o != nil {
				if err := o.Write(res); err != nil {
					fmt.Fprint(os.Stderr, "failed writing log\\n")
				}
			}
			if !success && exitOnError {
				cancel() // early exit - signal to stop searching for resources
				break
			}
		}

		for range validationResults { // allow resource finders to exit
		}

		result <- success
	}()

	return result
}

func kubeconform(cfg config.Config) int {
	var err error
	cpuProfileFile := os.Getenv("KUBECONFORM_CPUPROFILE_FILE")
	if cpuProfileFile != "" {
		f, err := os.Create(cpuProfileFile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		runtime.SetBlockProfileRate(1)

		defer pprof.StopCPUProfile()
	}

	useStdin := false
	if len(cfg.Files) == 0 || (len(cfg.Files) == 1 && cfg.Files[0] == "-") {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			log.Fatalf("failing to read data from stdin")
		}
		useStdin = true
	}

	var o output.Output
	if o, err = output.New(os.Stdout, cfg.OutputFormat, cfg.Summary, useStdin, cfg.Verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	var v validator.Validator
	v, err = validator.New(cfg.SchemaLocations, validator.Opts{
		Cache:                cfg.Cache,
		Debug:                cfg.Debug,
		SkipTLS:              cfg.SkipTLS,
		SkipKinds:            cfg.SkipKinds,
		RejectKinds:          cfg.RejectKinds,
		KubernetesVersion:    cfg.KubernetesVersion.String(),
		Strict:               cfg.Strict,
		IgnoreMissingSchemas: cfg.IgnoreMissingSchemas,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	validationResults := make(chan validator.Result)
	ctx, cancel := context.WithCancel(context.Background())
	successChan := processResults(cancel, o, validationResults, cfg.ExitOnError)

	var resourcesChan <-chan resource.Resource
	var errors <-chan error

	// Use the manifest with injected defaults for validation
	if cfg.InjectMissingDefaults {
		manifest, schema, err := loadManifestAndSchema(cfg.Files[0], cfg.SchemaLocations[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading manifest or schema: %s\\n", err)
			os.Exit(1)
		}

		// Inject defaults into the manifest
		injectDefaultsRecursively(schema, manifest)

		// Convert the modified manifest back to YAML for validation
		updatedManifest, err := yaml.Marshal(manifest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting updated manifest to YAML: %s\\n", err)
			os.Exit(1)
		}

		// Use a buffer as io.Reader to pass updated manifest
		manifestReader := bytes.NewReader(updatedManifest)

		// Use the updated manifest for validation
		resourcesChan, errors = resource.FromStream(ctx, "updatedManifest", manifestReader)
	} else {
		if useStdin {
			resourcesChan, errors = resource.FromStream(ctx, "stdin", os.Stdin)
		} else {
			resourcesChan, errors = resource.FromFiles(ctx, cfg.Files, cfg.IgnoreFilenamePatterns)
		}
	}

	// Process discovered resources across multiple workers
	wg := sync.WaitGroup{}
	for i := 0; i < cfg.NumberOfWorkers; i++ {
		wg.Add(1)
		go func(resources <-chan resource.Resource, validationResults chan<- validator.Result, v validator.Validator) {
			for res := range resources {
				validationResults <- v.ValidateResource(res)
			}
			wg.Done()
		}(resourcesChan, validationResults, v)
	}

	wg.Add(1)
	go func() {
		// Process errors while discovering resources
		for err := range errors {
			if err == nil {
				continue
			}

			if err, ok := err.(resource.DiscoveryError); ok {
				validationResults <- validator.Result{
					Resource: resource.Resource{Path: err.Path},
					Err:      err.Err,
					Status:   validator.Error,
				}
			} else {
				validationResults <- validator.Result{
					Resource: resource.Resource{},
					Err:      err,
					Status:   validator.Error,
				}
			}
			cancel()
		}
		wg.Done()
	}()

	wg.Wait()

	close(validationResults)
	success := <-successChan
	o.Flush()

	if !success {
		return 1
	}

	return 0
}

func main() {
	cfg, out, err := config.FromFlags(os.Args[0], os.Args[1:])
	if out != "" {
		o := os.Stderr
		errCode := 1
		if cfg.Help {
			o = os.Stdout
			errCode = 0
		}
		fmt.Fprintln(o, out)
		os.Exit(errCode)
	}

	if cfg.Version {
		fmt.Println(version)
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed parsing command line: %s\\n", err.Error())
		os.Exit(1)
	}

	// Inject defaults if the flag is enabled and validate the updated manifest
	os.Exit(kubeconform(cfg))
}
