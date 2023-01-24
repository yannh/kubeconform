// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"bytes"
	"io"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type doc struct {
	nodes []*yaml.Node
}

func NewDoc(nodes ...*yaml.Node) *doc {
	return &doc{nodes: nodes}
}

func ParseDoc(b []byte) (*doc, error) {
	br := bytes.NewReader(b)

	var nodes []*yaml.Node
	decoder := yaml.NewDecoder(br)
	for {
		node := &yaml.Node{}
		if err := decoder.Decode(node); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return &doc{nodes: nodes}, nil
}

func (d *doc) ToYAML() ([]byte, error) {
	var w bytes.Buffer
	encoder := yaml.NewEncoder(&w)
	for _, node := range d.nodes {
		if node.Kind == yaml.DocumentNode {
			if len(node.Content) == 0 {
				// These cause errors when we try to write them
				continue
			}
		}
		if err := encoder.Encode(node); err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}

func (d *doc) Elements() ([]*MapVariant, error) {
	return ExtractObjects(d.nodes...)
}
