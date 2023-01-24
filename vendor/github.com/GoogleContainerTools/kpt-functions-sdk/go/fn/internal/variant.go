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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type variantKind string

const (
	variantKindMap    variantKind = "Map"
	variantKindSlice  variantKind = "Slice"
	variantKindScalar variantKind = "Scalar"
)

type variant interface {
	GetKind() variantKind
	Node() *yaml.Node
}

// nodes are expected to be key1,value1,key2,value2,...
func buildMappingNode(nodes ...*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: nodes,
	}
}

func buildSequenceNode(nodes ...*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: nodes,
	}
}

func buildStringNode(s string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: s,
	}
}

func buildIntNode(i int) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int",
		Value: strconv.Itoa(i),
	}
}

func buildFloatNode(f float64) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!float",
		Value: strconv.FormatFloat(f, 'f', -1, 64),
	}
}

func buildBoolNode(b bool) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!bool",
		Value: strconv.FormatBool(b),
	}
}

func toVariant(n *yaml.Node) variant {
	switch n.Kind {
	case yaml.ScalarNode:
		return &scalarVariant{node: n}
	case yaml.MappingNode:
		return &MapVariant{node: n}
	case yaml.SequenceNode:
		return &sliceVariant{node: n}

	default:
		panic("unhandled yaml node kind")
	}
}

func ExtractObjects(nodes ...*yaml.Node) ([]*MapVariant, error) {
	var objects []*MapVariant

	for _, node := range nodes {
		switch node.Kind {
		case yaml.DocumentNode:
			children, err := ExtractObjects(node.Content...)
			if err != nil {
				return nil, err
			}
			objects = append(objects, children...)
		case yaml.MappingNode:
			objects = append(objects, &MapVariant{node: node})
		default:
			return nil, fmt.Errorf("unhandled node kind %v", node.Kind)
		}
	}
	return objects, nil
}

func TypedObjectToMapVariant(v interface{}) (*MapVariant, error) {
	// The built-in types only have json tags. We can't simply do ynode.Encode(v),
	// since it use the lowercased field name by default if no yaml tag is specified.
	// This affects both k8s built-in types (e.g. appsv1.Deployment) and any types
	// that depends on built-in types (e.g. metav1.ObjectMeta, corev1.PodTemplate).
	// To work around it, we rely on the json tags. We first convert v to
	// map[string]interface{} through json and then convert it to ynode.
	node, err := func() (*yaml.Node, error) {
		j, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var m map[string]interface{}
		if err = json.Unmarshal(j, &m); err != nil {
			return nil, err
		}

		node := &yaml.Node{}
		if err = node.Encode(m); err != nil {
			return nil, err
		}
		return node, nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to convert strong typed object to yaml node: %w", err)
	}

	mv := &MapVariant{node: node}
	mv.cleanupCreationTimestamp()
	err = mv.sortFields()
	return mv, err
}

func TypedObjectToSliceVariant(v interface{}) (*sliceVariant, error) {
	// The built-in types only have json tags. We can't simply do ynode.Encode(v),
	// since it use the lowercased field name by default if no yaml tag is specified.
	// This affects both k8s built-in types (e.g. appsv1.Deployment) and any types
	// that depends on built-in types (e.g. metav1.ObjectMeta, corev1.PodTemplate).
	// To work around it, we rely on the json tags. We first convert v to
	// map[string]interface{} through json and then convert it to ynode.
	node, err := func() (*yaml.Node, error) {
		j, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var l []interface{}
		if err = json.Unmarshal(j, &l); err != nil {
			return nil, err
		}

		node := &yaml.Node{}
		if err = node.Encode(l); err != nil {
			return nil, err
		}
		return node, nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to convert strong typed object to yaml node: %w", err)
	}

	return &sliceVariant{node: node}, nil
}

func MapVariantToTypedObject(mv *MapVariant, ptr interface{}) error {
	if ptr == nil || reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("ptr must be a pointer to an object")
	}

	// The built-in types only have json tags. We can't simply do mv.Node().Decode(ptr),
	// since it use the lowercased field name by default if no yaml tag is specified.
	// This affects both k8s built-in types (e.g. appsv1.Deployment) and any types
	// that depends on built-in types (e.g. metav1.ObjectMeta, corev1.PodTemplate).
	// To work around it, we rely on the json tags. We first convert mv to json
	// and then unmarshal it to ptr.
	j, err := yaml.NewRNode(mv.Node()).MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, ptr)
	return err
}
