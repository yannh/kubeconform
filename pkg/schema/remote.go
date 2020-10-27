package schema

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

const kubernetesJSONSchemaURLTmpl = "https://kubernetesjsonschema.dev/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json"

type remote string

func (r remote) Get(kind string, version string, kubernetesVersion string) (*Schema, error) {
	url, err := path(string(r), kind, version, kubernetesVersion, true)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no schema found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while downloading schema - received HTTP status %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed downloading schema at %s: %s", url, err)
	}

	schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(b))
	return (*Schema)(schema), err
}

// FromRemote TODO
func FromRemote(url string) Option {
	return func(r *Repository) {
		r.fetcher = append(r.fetcher, remote(url))
	}
}
