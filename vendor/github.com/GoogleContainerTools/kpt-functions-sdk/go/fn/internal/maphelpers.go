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

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (o *MapVariant) GetRNode(fields ...string) (*yaml.RNode, bool, error) {
	rn := &yaml.RNode{}
	val, found, err := o.GetNestedValue(fields...)
	if err != nil || !found {
		return nil, found, err
	}
	rn.SetYNode(val.Node())
	return rn, found, err
}

func (o *MapVariant) GetNestedValue(fields ...string) (variant, bool, error) {
	current := o
	n := len(fields)
	for i := 0; i < n; i++ {
		entry, found := current.getVariant(fields[i])
		if !found {
			return nil, found, nil
		}

		if i == n-1 {
			return entry, true, nil
		}
		entryM, ok := entry.(*MapVariant)
		if !ok {
			return nil, found, fmt.Errorf("wrong type, got: %T", entry)
		}
		current = entryM
	}
	return nil, false, fmt.Errorf("unexpected code reached")
}

func (o *MapVariant) SetNestedValue(val variant, fields ...string) error {
	current := o
	n := len(fields)
	var err error
	for i := 0; i < n; i++ {
		if i == n-1 {
			current.set(fields[i], val)
		} else {
			current, _, err = current.getMap(fields[i], true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *MapVariant) GetNestedMap(fields ...string) (*MapVariant, bool, error) {
	v, found, err := o.GetNestedValue(fields...)
	if err != nil || !found {
		return nil, found, err
	}
	mv, ok := v.(*MapVariant)
	if !ok {
		return nil, found, fmt.Errorf("wrong type, got: %T", v)
	}
	return mv, found, err
}

func (o *MapVariant) SetNestedMap(m *MapVariant, fields ...string) error {
	return o.SetNestedValue(m, fields...)
}

func (o *MapVariant) GetNestedStringMap(fields ...string) (map[string]string, bool, error) {
	v, found, err := o.GetNestedValue(fields...)
	if err != nil || !found {
		return nil, found, err
	}
	children := v.Node().Content
	if len(children)%2 != 0 {
		return nil, found, fmt.Errorf("invalid yaml map node")
	}
	m := make(map[string]string, len(children)/2)
	for i := 0; i < len(children); i = i + 2 {
		m[children[i].Value] = children[i+1].Value
	}
	return m, found, nil
}

func (o *MapVariant) SetNestedStringMap(m map[string]string, fields ...string) error {
	return o.SetNestedMap(NewStringMapVariant(m), fields...)
}

func (o *MapVariant) GetNestedScalar(fields ...string) (*scalarVariant, bool, error) {
	node, found, err := o.GetNestedValue(fields...)
	if err != nil || !found {
		return nil, found, err
	}
	nodeS, ok := node.(*scalarVariant)
	if !ok {
		return nil, found, fmt.Errorf("incorrect type, was %T", node)
	}
	return nodeS, found, nil
}

func (o *MapVariant) GetNestedString(fields ...string) (string, bool, error) {
	scalar, found, err := o.GetNestedScalar(fields...)
	if err != nil || !found {
		return "", found, err
	}
	sv, isString := scalar.StringValue()
	if isString {
		return sv, found, nil
	}
	return "", found, fmt.Errorf("node was not a string, was %v", scalar.node.Tag)
}

func (o *MapVariant) SetNestedString(s string, fields ...string) error {
	return o.SetNestedValue(newStringScalarVariant(s), fields...)
}

func (o *MapVariant) GetNestedBool(fields ...string) (bool, bool, error) {
	scalar, found, err := o.GetNestedScalar(fields...)
	if err != nil || !found {
		return false, found, err
	}
	bv, isBool := scalar.BoolValue()
	if isBool {
		return bv, found, nil
	}
	return false, found, fmt.Errorf("node was not a bool, was %v", scalar.Node().Tag)
}

func (o *MapVariant) SetNestedBool(b bool, fields ...string) error {
	return o.SetNestedValue(newBoolScalarVariant(b), fields...)
}

func (o *MapVariant) GetNestedInt(fields ...string) (int, bool, error) {
	scalar, found, err := o.GetNestedScalar(fields...)
	if err != nil || !found {
		return 0, found, err
	}
	iv, isInt := scalar.IntValue()
	if isInt {
		return iv, found, nil
	}
	return 0, found, fmt.Errorf("node was not a int, was %v", scalar.node.Tag)
}

func (o *MapVariant) SetNestedInt(i int, fields ...string) error {
	return o.SetNestedValue(newIntScalarVariant(i), fields...)
}

func (o *MapVariant) GetNestedFloat(fields ...string) (float64, bool, error) {
	scalar, found, err := o.GetNestedScalar(fields...)
	if err != nil || !found {
		return 0, found, err
	}
	fv, isFloat := scalar.FloatValue()
	if isFloat {
		return fv, found, nil
	}
	return 0, found, fmt.Errorf("node was not a float, was %v", scalar.node.Tag)
}

func (o *MapVariant) SetNestedFloat(f float64, fields ...string) error {
	return o.SetNestedValue(newFloatScalarVariant(f), fields...)
}

func (o *MapVariant) GetNestedSlice(fields ...string) (*sliceVariant, bool, error) {
	node, found, err := o.GetNestedValue(fields...)
	if err != nil || !found {
		return nil, found, err
	}
	nodeS, ok := node.(*sliceVariant)
	if !ok {
		return nil, found, fmt.Errorf("incorrect type, was %T", node)
	}
	return nodeS, found, err
}

func (o *MapVariant) SetNestedSlice(s *sliceVariant, fields ...string) error {
	return o.SetNestedValue(s, fields...)
}

func (o *MapVariant) RemoveNestedField(fields ...string) (bool, error) {
	current := o
	n := len(fields)
	for i := 0; i < n; i++ {
		entry, found := current.getVariant(fields[i])
		if !found {
			return false, nil
		}

		if i == n-1 {
			return current.remove(fields[i])
		}
		switch entry := entry.(type) {
		case *MapVariant:
			current = entry
		default:
			return false, fmt.Errorf("value is of unexpected type %T", entry)
		}
	}
	return false, fmt.Errorf("unexpected code reached")
}

func (o *MapVariant) getMap(field string, create bool) (*MapVariant, bool, error) {
	node, found := o.getVariant(field)

	if !found {
		if !create {
			return nil, found, nil
		}
		keyNode := buildStringNode(field)
		valueNode := buildMappingNode()
		o.node.Content = append(o.node.Content, keyNode, valueNode)
		valueVariant := &MapVariant{node: valueNode}
		return valueVariant, found, nil
	}

	if node, ok := node.(*MapVariant); ok {
		return node, found, nil
	}
	return nil, found, fmt.Errorf("incorrect type, was %T", node)
}
