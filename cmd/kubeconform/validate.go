package kubeconform

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

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
					log.Fatal("failed writing log: ", err)
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

func Validate(cfg config.Config, out string) error {

	if out != "" {
		if cfg.Help {
			fmt.Fprintln(cfg.Stream.Output, out)
			return nil
		}
		return fmt.Errorf("config out is not empty")
	}

	if cfg.Version {
		fmt.Fprintln(cfg.Stream.Output, out)
		return nil
	}

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

	o, err := output.New(cfg.Stream.Output, cfg.OutputFormat, cfg.Summary, useStdin, cfg.Verbose)
	if err != nil {
		return fmt.Errorf("failed to get output: %s", err.Error())
	}
	v, err := validator.New(cfg.SchemaLocations, validator.Opts{
		Cache:                cfg.Cache,
		Debug:                cfg.Debug,
		SkipTLS:              cfg.SkipTLS,
		SkipKinds:            cfg.SkipKinds,
		RejectKinds:          cfg.RejectKinds,
		KubernetesVersion:    cfg.KubernetesVersion,
		Strict:               cfg.Strict,
		IgnoreMissingSchemas: cfg.IgnoreMissingSchemas,
	})
	if err != nil {
		return fmt.Errorf("failed to validate: %s", err.Error())
	}

	validationResults := make(chan validator.Result)
	ctx, cancel := context.WithCancel(context.Background())
	successChan := processResults(cancel, o, validationResults, cfg.ExitOnError)

	var resourcesChan <-chan resource.Resource
	var errors <-chan error
	if useStdin {
		resourcesChan, errors = resource.FromStream(ctx, "stdin", cfg.Stream.Input)
	} else {
		resourcesChan, errors = resource.FromFiles(ctx, cfg.Files, cfg.IgnoreFilenamePatterns)
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
		return fmt.Errorf("failed to process results")
	}

	return nil
}
