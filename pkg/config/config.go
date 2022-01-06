package config

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Cache                  string
	CPUProfileFile         string
	ExitOnError            bool
	Files                  []string
	SchemaLocations        []string
	SkipTLS                bool
	SkipKinds              map[string]struct{}
	RejectKinds            map[string]struct{}
	OutputFormat           string
	KubernetesVersion      string
	NumberOfWorkers        int
	Summary                bool
	Strict                 bool
	Verbose                bool
	IgnoreMissingSchemas   bool
	IgnoreFilenamePatterns []string
	Help                   bool
	Version                bool
}

type arrayParam []string

func (ap *arrayParam) String() string {
	return strings.Join(*ap, " - ")
}

func (ap *arrayParam) Set(value string) error {
	*ap = append(*ap, value)
	return nil
}

func splitCSV(csvStr string) map[string]struct{} {
	splitValues := strings.Split(csvStr, ",")
	valuesMap := map[string]struct{}{}

	for _, kind := range splitValues {
		if len(kind) > 0 {
			valuesMap[kind] = struct{}{}
		}
	}

	return valuesMap
}

// FromFlags retrieves kubeconform's runtime configuration from the command-line parameters
func FromFlags(progName string, args []string) (Config, string, error) {
	var schemaLocationsParam, ignoreFilenamePatterns arrayParam
	var skipKindsCSV, rejectKindsCSV string
	flags := flag.NewFlagSet(progName, flag.ExitOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	c := Config{}
	c.Files = []string{}

	flags.StringVar(&c.KubernetesVersion, "kubernetes-version", "master", "version of Kubernetes to validate against, e.g.: 1.18.0")
	flags.Var(&schemaLocationsParam, "schema-location", "override schemas location search path (can be specified multiple times)")
	flags.StringVar(&skipKindsCSV, "skip", "", "comma-separated list of kinds to ignore")
	flags.StringVar(&rejectKindsCSV, "reject", "", "comma-separated list of kinds to reject")
	flags.BoolVar(&c.ExitOnError, "exit-on-error", false, "immediately stop execution when the first error is encountered")
	flags.BoolVar(&c.IgnoreMissingSchemas, "ignore-missing-schemas", false, "skip files with missing schemas instead of failing")
	flags.Var(&ignoreFilenamePatterns, "ignore-filename-pattern", "regular expression specifying paths to ignore (can be specified multiple times)")
	flags.BoolVar(&c.Summary, "summary", false, "print a summary at the end (ignored for junit output)")
	flags.IntVar(&c.NumberOfWorkers, "n", 4, "number of goroutines to run concurrently")
	flags.BoolVar(&c.Strict, "strict", false, "disallow additional properties not in schema")
	flags.StringVar(&c.OutputFormat, "output", "text", "output format - json, junit, tap, text")
	flags.BoolVar(&c.Verbose, "verbose", false, "print results for all resources (ignored for tap and junit output)")
	flags.BoolVar(&c.SkipTLS, "insecure-skip-tls-verify", false, "disable verification of the server's SSL certificate. This will make your HTTPS connections insecure")
	flags.StringVar(&c.Cache, "cache", "", "cache schemas downloaded via HTTP to this folder")
	flags.StringVar(&c.CPUProfileFile, "cpu-prof", "", "debug - log CPU profiling to file")
	flags.BoolVar(&c.Help, "h", false, "show help information")
	flags.BoolVar(&c.Version, "v", false, "show version information")
	flags.Usage = func() {
		fmt.Fprintf(&buf, "Usage: %s [OPTION]... [FILE OR FOLDER]...\n", progName)
		flags.PrintDefaults()
	}

	err := flags.Parse(args)

	c.SkipKinds = splitCSV(skipKindsCSV)
	c.RejectKinds = splitCSV(rejectKindsCSV)
	c.IgnoreFilenamePatterns = ignoreFilenamePatterns
	c.SchemaLocations = schemaLocationsParam
	c.Files = flags.Args()

	if c.Help {
		flags.Usage()
	}

	return c, buf.String(), err
}
