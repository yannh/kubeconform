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

package fn

import (
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// byteReadWriter wraps kio.ByteReadWriter
type byteReadWriter struct {
	kio.ByteReadWriter
}

// Read decodes input bytes into a ResourceList
func (rw *byteReadWriter) Read() (*ResourceList, error) {
	nodes, err := rw.ByteReadWriter.Read()
	if err != nil {
		return nil, err
	}
	var items KubeObjects
	for _, n := range nodes {
		obj, err := ParseKubeObject([]byte(n.MustString()))
		if err != nil {
			return nil, err
		}
		items = append(items, obj)
	}
	obj, err := ParseKubeObject([]byte(rw.ByteReadWriter.FunctionConfig.MustString()))
	if err != nil {
		return nil, err
	}
	return &ResourceList{
		Items:          items,
		FunctionConfig: obj,
	}, nil
}

// Write writes a ResourceList into bytes
func (rw *byteReadWriter) Write(rl *ResourceList) error {
	if len(rl.Results) > 0 {
		b, err := yaml.Marshal(rl.Results)
		if err != nil {
			return errors.Wrap(err)
		}
		y, err := yaml.Parse(string(b))
		if err != nil {
			return errors.Wrap(err)
		}
		rw.Results = y
	}
	var nodes []*yaml.RNode
	for _, item := range rl.Items {
		node, err := yaml.Parse(item.String())
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
	}
	return rw.ByteReadWriter.Write(nodes)
}
