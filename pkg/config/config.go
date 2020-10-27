package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config TODO
type Config struct {
	RootPaths            []string
	SchemaLocations      []string
	SkipKinds            map[string]struct{}
	OutputFormat         string
	KubernetesVersion    string
	NumberOfWorkers      int
	Summary              bool
	Strict               bool
	Verbose              bool
	IgnoreMissingSchemas bool
	Help                 bool
}

type arrayParam []string

func (ap *arrayParam) String() string {
	return strings.Join(*ap, " - ")
}

func (ap *arrayParam) Set(value string) error {
	*ap = append(*ap, value)
	return nil
}

func skipKinds(skipKindsCSV string) map[string]struct{} {
	splitKinds := strings.Split(skipKindsCSV, ",")
	skipKinds := map[string]struct{}{}

	for _, kind := range splitKinds {
		if len(kind) > 0 {
			skipKinds[kind] = struct{}{}
		}
	}

	return skipKinds
}

// FromFlags TODO
func FromFlags() Config {
	c := Config{}

	flag.StringVar(&c.KubernetesVersion, "kubernetes-version", "1.18.0", "version of Kubernetes to validate against")

	var schemaLocationsParam arrayParam
	flag.Var(&schemaLocationsParam, "schema-location", "override schemas location search path (can be specified multiple times)")
	c.SchemaLocations = []string(schemaLocationsParam)

	var skipKindsCSV string
	flag.StringVar(&skipKindsCSV, "skip", "", "comma-separated list of kinds to ignore")
	c.SkipKinds = skipKinds(skipKindsCSV)

	flag.BoolVar(&c.IgnoreMissingSchemas, "ignore-missing-schemas", false, "skip files with missing schemas instead of failing")
	flag.BoolVar(&c.Summary, "summary", false, "print a summary at the end")
	flag.IntVar(&c.NumberOfWorkers, "n", -1, "number of goroutines to run concurrently")

	flag.BoolVar(&c.Strict, "strict", false, "disallow additional properties not in schema")
	flag.StringVar(&c.OutputFormat, "output", "text", "output format - text, json")
	flag.BoolVar(&c.Verbose, "verbose", false, "print results for all resources")
	flag.BoolVar(&c.Help, "h", false, "show help information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [FILE OR FOLDER]...\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if c.Help {
		flag.Usage()
	}

	return c
}
