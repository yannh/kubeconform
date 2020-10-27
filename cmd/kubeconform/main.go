package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/output"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/schema"
	"github.com/yannh/kubeconform/pkg/validation"
)

func processResults(o output.Output, results <-chan validation.Result) <-chan bool {
	done := make(chan bool)

	go func() {
		ret := true
		for result := range results {
			sig, err := result.Signature()
			if err != nil {
				// fmt.Fprint(os.Stderr, "failed writing log\n")
			}
			if result.Status == validation.Invalid {
				ret = false
			}
			err = o.Write(result.Path, sig.Kind, sig.Metadata.Name, sig.APIVersion, result.Err, result.Status == validation.Skipped)
			if err != nil {
				fmt.Fprint(os.Stderr, "failed writing log\n")
			}
		}
		done <- ret
	}()

	return done
}

func main() {
	conf := config.FromFlags()
	if conf.Help {
		os.Exit(0)
	}

	if conf.NumberOfWorkers == -1 {
		conf.NumberOfWorkers = runtime.NumCPU()
	}

	var opts []schema.Option
	for _, loc := range conf.SchemaLocations {
		if strings.HasPrefix(loc, "http://") || strings.HasPrefix(loc, "https://") {
			opts = append(opts, schema.FromRemote(loc))
		} else {
			opts = append(opts, schema.FromFS(loc))
		}
	}

	schemas := schema.New(opts...)
	validator := validation.New(schemas, conf)
	resources := resource.Discover(flag.Args()...)

	results := make(chan validation.Result)
	var wg sync.WaitGroup
	wg.Add(conf.NumberOfWorkers)
	for i := 0; i < conf.NumberOfWorkers; i++ {
		go func() {
			for batch := range resources {
				for _, res := range batch {
					results <- validator.Validate(res)
				}
			}
			wg.Done()
		}()
	}

	output, err := output.New(conf.OutputFormat, conf.Summary, conf.Verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting output: %s\n", err)
		os.Exit(1)
	}

	done := processResults(output, results)

	wg.Wait()
	close(results)
	if <-done {
		os.Exit(0)
	}

	os.Exit(1)
}
