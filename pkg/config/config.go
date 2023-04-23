package config

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Stream struct {
	Input  io.Reader
	Output io.Writer
	Error  io.Writer
}

type Config struct {
	Cache                  string   `yaml:"cache" json:"cache"`
	Debug                  bool     `yaml:"debug" json:"debug"`
	ExitOnError            bool     `yaml:"exitOnError" json:"exitOnError"`
	Files                  []string `yaml:"files" json:"files"`
	Help                   bool     `yaml:"help" json:"help"`
	IgnoreFilenamePatterns []string `yaml:"ignoreFilenamePatterns" json:"ignoreFilenamePatterns"`
	IgnoreMissingSchemas   bool     `yaml:"ignoreMissingSchemas" json:"ignoreMissingSchemas"`
	KubernetesVersion      string   `yaml:"kubernetesVersion" json:"kubernetesVersion"`
	NumberOfWorkers        int      `yaml:"numberOfWorkers" json:"numberOfWorkers"`
	OutputFormat           string   `yaml:"output" json:"output"`
	RejectKinds            map[string]struct{}
	RejectKindsNG          []string `yaml:"reject" json:"reject"`
	SchemaLocations        []string `yaml:"schemaLocations" json:"schemaLocations"`
	SkipKinds              map[string]struct{}
	SkipKindsNG            []string `yaml:"skip" json:"skip"`
	SkipTLS                bool     `yaml:"insecureSkipTLSVerify" json:"insecureSkipTLSVerify"`
	Strict                 bool     `yaml:"strict" json:"strict"`
	Summary                bool     `yaml:"summary" json:"summary"`
	Verbose                bool     `yaml:"verbose" json:"verbose"`
	Version                bool     `yaml:"version" json:"version"`
	Stream                 *Stream
}

// LoadNGConfig allows to introduce new config format/structure while keeping backward compatibility.
func (c *Config) LoadNGConfig() error {
	// SkipKindsNG.
	loadKinds(c.SkipKinds, c.SkipKindsNG)

	// RejectKindsNG.
	loadKinds(c.RejectKinds, c.RejectKindsNG)

	return nil
}

func loadKinds(dest map[string]struct{}, src []string) {
	for _, kind := range src {
		if len(kind) > 0 {
			dest[kind] = struct{}{}
		}
	}
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
	valuesMap := map[string]struct{}{}
	loadKinds(valuesMap, strings.Split(csvStr, ","))
	return valuesMap
}

// FromFlags retrieves kubeconform's runtime configuration from the command-line parameters
func FromFlags(progName string, args []string) (Config, string, error) {
	var schemaLocationsParam, ignoreFilenamePatterns arrayParam
	var skipKindsCSV, rejectKindsCSV string
	flags := flag.NewFlagSet(progName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	c := Config{}
	c.Files = []string{}
	c.Stream = &Stream{os.Stdin, os.Stdout, os.Stderr}

	flags.SetOutput(c.Stream.Output)
	flags.StringVar(&c.KubernetesVersion, "kubernetes-version", "master", "version of Kubernetes to validate against, e.g.: 1.18.0")
	flags.Var(&schemaLocationsParam, "schema-location", "override schemas location search path (can be specified multiple times)")
	flags.StringVar(&skipKindsCSV, "skip", "", "comma-separated list of kinds or GVKs to ignore")
	flags.StringVar(&rejectKindsCSV, "reject", "", "comma-separated list of kinds or GVKs to reject")
	flags.BoolVar(&c.Debug, "debug", false, "print debug information")
	flags.BoolVar(&c.ExitOnError, "exit-on-error", false, "immediately stop execution when the first error is encountered")
	flags.BoolVar(&c.IgnoreMissingSchemas, "ignore-missing-schemas", false, "skip files with missing schemas instead of failing")
	flags.Var(&ignoreFilenamePatterns, "ignore-filename-pattern", "regular expression specifying paths to ignore (can be specified multiple times)")
	flags.BoolVar(&c.Summary, "summary", false, "print a summary at the end (ignored for junit output)")
	flags.IntVar(&c.NumberOfWorkers, "n", 4, "number of goroutines to run concurrently")
	flags.BoolVar(&c.Strict, "strict", false, "disallow additional properties not in schema or duplicated keys")
	flags.StringVar(&c.OutputFormat, "output", "text", "output format - json, junit, pretty, tap, text")
	flags.BoolVar(&c.Verbose, "verbose", false, "print results for all resources (ignored for tap and junit output)")
	flags.BoolVar(&c.SkipTLS, "insecure-skip-tls-verify", false, "disable verification of the server's SSL certificate. This will make your HTTPS connections insecure")
	flags.StringVar(&c.Cache, "cache", "", "cache schemas downloaded via HTTP to this folder")
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
