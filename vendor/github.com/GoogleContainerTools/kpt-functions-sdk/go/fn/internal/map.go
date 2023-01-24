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
	"fmt"
	"log"
	"sort"

	"k8s.io/klog/v2"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func NewMap(node *yaml.Node) *MapVariant {
	if node == nil {
		node = &yaml.Node{
			Kind: yaml.MappingNode,
		}
	}
	return &MapVariant{node: node}
}

func NewStringMapVariant(m map[string]string) *MapVariant {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
	}
	for k, v := range m {
		node.Content = append(node.Content, buildStringNode(k), buildStringNode(v))
	}
	return &MapVariant{node: node}
}

type MapVariant struct {
	node *yaml.Node
}

func (o *MapVariant) GetKind() variantKind {
	return variantKindMap
}

func (o *MapVariant) Node() *yaml.Node {
	return o.node
}

func (o *MapVariant) Entries() (map[string]variant, error) {
	entries := make(map[string]variant)

	ynode := o.node
	children := ynode.Content
	if len(children)%2 != 0 {
		return nil, fmt.Errorf("unexpected number of children for map %d", len(children))
	}

	for i := 0; i < len(children); i += 2 {
		keyNode := children[i]
		valueNode := children[i+1]

		keyVariant := toVariant(keyNode)
		valueVariant := toVariant(valueNode)

		switch keyVariant := keyVariant.(type) {
		case *scalarVariant:
			sv, isString := keyVariant.StringValue()
			if isString {
				entries[sv] = valueVariant
			} else {
				return nil, fmt.Errorf("key was not a string %v", keyVariant)
			}
		default:
			return nil, fmt.Errorf("unexpected variant kind %T", keyVariant)
		}
	}
	return entries, nil
}

func asString(node *yaml.Node) (string, bool) {
	if node.Kind == yaml.ScalarNode && (node.Tag == "!!str" || node.Tag == "") {
		return node.Value, true
	}
	return "", false
}

func (o *MapVariant) getVariant(key string) (variant, bool) {
	valueNode, found := getValueNode(o.node, key)
	if !found {
		return nil, found
	}

	v := toVariant(valueNode)
	return v, true
}

func getValueNode(m *yaml.Node, key string) (*yaml.Node, bool) {
	children := m.Content
	if len(children)%2 != 0 {
		log.Fatalf("unexpected number of children for map %d", len(children))
	}

	for i := 0; i < len(children); i += 2 {
		keyNode := children[i]

		k, ok := asString(keyNode)
		if ok && k == key {
			valueNode := children[i+1]
			return valueNode, true
		}
	}
	return nil, false
}

func (o *MapVariant) set(key string, val variant) {
	o.setYAMLNode(key, val.Node())
}

func (o *MapVariant) setYAMLNode(key string, node *yaml.Node) {
	children := o.node.Content
	if len(children)%2 != 0 {
		log.Fatalf("unexpected number of children for map %d", len(children))
	}

	for i := 0; i < len(children); i += 2 {
		keyNode := children[i]

		k, ok := asString(keyNode)
		if ok && k == key {
			// TODO: Copy comments?
			oldNode := children[i+1]
			children[i+1] = node
			children[i+1].FootComment = oldNode.FootComment
			children[i+1].HeadComment = oldNode.HeadComment
			children[i+1].LineComment = oldNode.LineComment
			return
		}
	}

	o.node.Content = append(o.node.Content, buildStringNode(key), node)
}

func (o *MapVariant) remove(key string) (bool, error) {
	removed := false

	children := o.node.Content
	if len(children)%2 != 0 {
		return false, fmt.Errorf("unexpected number of children for map %d", len(children))
	}

	var keep []*yaml.Node
	for i := 0; i < len(children); i += 2 {
		keyNode := children[i]

		k, ok := asString(keyNode)
		if ok && k == key {
			removed = true
			continue
		}

		keep = append(keep, children[i], children[i+1])
	}

	o.node.Content = keep

	return removed, nil
}

// remove field metadata.creationTimestamp when it's null.
func (o *MapVariant) cleanupCreationTimestamp() {
	if o.node.Kind != yaml.MappingNode {
		return
	}
	scalar, found, err := o.GetNestedScalar("metadata", "creationTimestamp")
	if err != nil || !found {
		return
	}
	if scalar.IsNull() {
		_, _ = o.RemoveNestedField("metadata", "creationTimestamp")
	}
}

// sortFields tried to sort fields that it understands. e.g. data should come
// after apiVersion, kind and metadata in corev1.ConfigMap.
func (o *MapVariant) sortFields() error {
	return sortFields(o.node)
}

func sortFields(ynode *yaml.Node) error {
	if ynode.Kind == yaml.SequenceNode {
		for _, child := range ynode.Content {
			if err := sortFields(child); err != nil {
				return err
			}
		}
		return nil
	}
	if ynode.Kind != yaml.MappingNode {
		return nil
	}

	pairs, err := ynodeToYamlKeyValuePairs(ynode)
	if err != nil {
		return fmt.Errorf("unable to sort fields in yaml: %w", err)
	}
	for _, pair := range pairs {
		if err = sortFields(pair.value); err != nil {
			return err
		}
	}
	sort.Sort(pairs)
	ynode.Content = yamlKeyValuePairsToYnode(pairs)
	return nil
}

func ynodeToYamlKeyValuePairs(ynode *yaml.Node) (yamlKeyValuePairs, error) {
	if len(ynode.Content)%2 != 0 {
		return nil, fmt.Errorf("invalid number of nodes: %d", len(ynode.Content))
	}

	var pairs yamlKeyValuePairs
	for i := 0; i < len(ynode.Content); i += 2 {
		pairs = append(pairs, &yamlKeyValuePair{name: ynode.Content[i], value: ynode.Content[i+1]})
	}
	return pairs, nil
}

func yamlKeyValuePairsToYnode(pairs yamlKeyValuePairs) []*yaml.Node {
	var nodes []*yaml.Node
	for _, pair := range pairs {
		nodes = append(nodes, pair.name, pair.value)
	}
	return nodes
}

type yamlKeyValuePair struct {
	name  *yaml.Node
	value *yaml.Node
}

type yamlKeyValuePairs []*yamlKeyValuePair

func (nodes yamlKeyValuePairs) Len() int { return len(nodes) }

func (nodes yamlKeyValuePairs) Less(i, j int) bool {
	iIndex, iFound := yaml.FieldOrder[nodes[i].name.Value]
	jIndex, jFound := yaml.FieldOrder[nodes[j].name.Value]
	if iFound && jFound {
		return iIndex < jIndex
	}
	if iFound {
		return true
	}
	if jFound {
		return false
	}

	if nodes[i].name != nodes[j].name {
		return nodes[i].name.Value < nodes[j].name.Value
	}
	return false
}

func (nodes yamlKeyValuePairs) Swap(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] }

// UpsertMap will return the field as a map if it exists and is a map,
// otherwise it will insert a map at the specified field.
// Note that if the value exists but is not a map, it will be replaced with a map.
func (o *MapVariant) UpsertMap(field string) *MapVariant {
	m := o.GetMap(field)
	if m != nil {
		return m
	}

	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: field,
		Tag:   "!!str",
	}
	valueNode := &yaml.Node{
		Kind: yaml.MappingNode,
	}
	o.node.Content = append(o.node.Content, keyNode, valueNode)
	return &MapVariant{node: valueNode}
}

// GetMap will return the field as a map if it exists and is a map,
// otherwise it will return nil.
// Note that if the value exists but is not a map, nil will be returned.
func (o *MapVariant) GetMap(field string) *MapVariant {
	node, found := o.getVariant(field)

	if found {
		switch node := node.(type) {
		case *MapVariant:
			return node

		default:
			klog.Warningf("getting value of unexpected type, got %T, want map", node)
		}
	}

	return nil
}
