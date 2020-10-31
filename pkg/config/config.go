package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Files                []string
	SchemaLocations      []string
	SkipKinds            map[string]bool
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

func skipKinds(skipKindsCSV string) map[string]bool {
	splitKinds := strings.Split(skipKindsCSV, ",")
	skipKinds := map[string]bool{}

	for _, kind := range splitKinds {
		if len(kind) > 0 {
			skipKinds[kind] = true
		}
	}

	return skipKinds
}

func FromFlags() Config {
	var schemaLocationsParam arrayParam
	var skipKindsCSV string

	c := Config{}
	c.Files = []string{}

	flag.StringVar(&c.KubernetesVersion, "kubernetes-version", "1.18.0", "version of Kubernetes to validate against")
	flag.Var(&schemaLocationsParam, "schema-location", "override schemas location search path (can be specified multiple times)")
	flag.StringVar(&skipKindsCSV, "skip", "", "comma-separated list of kinds to ignore")
	flag.BoolVar(&c.IgnoreMissingSchemas, "ignore-missing-schemas", false, "skip files with missing schemas instead of failing")
	flag.BoolVar(&c.Summary, "summary", false, "print a summary at the end")
	flag.IntVar(&c.NumberOfWorkers, "n", 4, "number of goroutines to run concurrently")
	flag.BoolVar(&c.Strict, "strict", false, "disallow additional properties not in schema")
	flag.StringVar(&c.OutputFormat, "output", "text", "output format - text, json")
	flag.BoolVar(&c.Verbose, "verbose", false, "print results for all resources")
	flag.BoolVar(&c.Help, "h", false, "show help information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [FILE OR FOLDER]...\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	c.SkipKinds = skipKinds(skipKindsCSV)
	c.SchemaLocations = schemaLocationsParam
	if len(c.SchemaLocations) == 0 {
		c.SchemaLocations = append(c.SchemaLocations, "https://kubernetesjsonschema.dev") // if not specified, default behaviour is to use kubernetesjson-schema.dev as registry
	}

	for _, file := range flag.Args() {
		c.Files = append(c.Files, file)
	}

	if c.Help {
		flag.Usage()
	}

	return c
}
