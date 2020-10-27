package schema

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/xeipuuv/gojsonschema"
)

type fs string

func path(tpl, resourceKind, resourceAPIVersion, k8sVersion string, strict bool) (string, error) {
	normalisedVersion := k8sVersion
	if normalisedVersion != "master" {
		normalisedVersion = "v" + normalisedVersion
	}

	strictSuffix := ""
	if strict {
		strictSuffix = "-strict"
	}

	groupParts := strings.Split(resourceAPIVersion, "/")
	versionParts := strings.Split(groupParts[0], ".")

	kindSuffix := "-" + strings.ToLower(versionParts[0])
	if len(groupParts) > 1 {
		kindSuffix += "-" + strings.ToLower(groupParts[1])
	}

	tmpl, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", err
	}

	tplData := struct {
		NormalizedVersion string
		StrictSuffix      string
		ResourceKind      string
		KindSuffix        string
	}{
		normalisedVersion,
		strictSuffix,
		strings.ToLower(resourceKind),
		kindSuffix,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tplData)

	return buf.String(), err
}

func (f fs) Get(kind string, version string, kubernetesVersion string) (*Schema, error) {
	p, err := path(string(f), kind, version, kubernetesVersion, true)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(b))

	return (*Schema)(schema), err
}

// FromFS TODO
func FromFS(path string) Option {
	return func(r *Repository) {
		r.fetcher = append(r.fetcher, fs(path))
	}
}
