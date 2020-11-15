package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

func processResults(ctx context.Context, o output.Output, validationResults <-chan validator.Result, exitOnError bool) <-chan bool {
	success := true
	result := make(chan bool)

	go func() {
		for res := range validationResults {
			if res.Status == validator.Error || res.Status == validator.Invalid {
				success = false
			}
			if o != nil {
				if err := o.Write(res); err != nil {
					fmt.Fprint(os.Stderr, "failed writing log\n")
				}
			}
			if success == false && exitOnError {
				ctx.Done() // early exit - signal to stop searching for resources
				break
			}
		}

		for range validationResults { // allow resource finders to exit
		}

		result <- success
		close(result)
	}()

	return result
}

func realMain() int {
	cfg, out, err := config.FromFlags(os.Args[0], os.Args[1:])
	if out != "" {
		fmt.Println(out)
		return 1
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "failed parsing command line: %s\n", err.Error())
		return 1
	}

	// Detect whether we have data being piped through stdin
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(cfg.Files) == 1 && cfg.Files[0] == "-" {
		isStdin = true
	}

	var o output.Output
	if o, err = output.New(cfg.OutputFormat, cfg.Summary, isStdin, cfg.Verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	v, err := validator.New(cfg.SchemaLocations, validator.Opts{
		SkipTLS:              cfg.SkipTLS,
		SkipKinds:            cfg.SkipKinds,
		RejectKinds:          cfg.RejectKinds,
		KubernetesVersion:    cfg.KubernetesVersion,
		Strict:               cfg.Strict,
		IgnoreMissingSchemas: cfg.IgnoreMissingSchemas,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	validationResults := make(chan validator.Result)
	ctx := context.Background()
	successChan := processResults(ctx, o, validationResults, cfg.ExitOnError)

	var resourcesChan <-chan resource.Resource
	var errors <-chan error
	if isStdin {
		resourcesChan, errors = resource.FromStream(ctx, "stdin", os.Stdin)
	} else {
		resourcesChan, errors = resource.FromFiles(ctx, cfg.IgnoreFilenamePatterns, cfg.Files...)
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
			ctx.Done()
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
	os.Exit(realMain())
}
