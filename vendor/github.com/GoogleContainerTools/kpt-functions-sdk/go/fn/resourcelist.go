// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package fn

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn/internal"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// ResourceList is a Kubernetes list type used as the primary data interchange format
// in the Configuration Functions Specification:
// https://github.com/kubernetes-sigs/kustomize/blob/master/cmd/config/docs/api-conventions/functions-spec.md
// This framework facilitates building functions that receive and emit ResourceLists,
// as required by the specification.
type ResourceList struct {
	// Items is the ResourceList.items input and output value.
	//
	// e.g. given the function input:
	//
	//    kind: ResourceList
	//    items:
	//    - kind: Deployment
	//      ...
	//    - kind: Service
	//      ...
	//
	// Items will be a slice containing the Deployment and Service resources
	// Mutating functions will alter this field during processing.
	// This field is required.
	Items KubeObjects `yaml:"items" json:"items"`

	// FunctionConfig is the ResourceList.functionConfig input value.
	//
	// e.g. given the input:
	//
	//    kind: ResourceList
	//    functionConfig:
	//      kind: Example
	//      spec:
	//        foo: var
	//
	// FunctionConfig will contain the RNodes for the Example:
	//      kind: Example
	//      spec:
	//        foo: var
	FunctionConfig *KubeObject `yaml:"functionConfig,omitempty" json:"functionConfig,omitempty"`

	// Results is ResourceList.results output value.
	// Validating functions can optionally use this field to communicate structured
	// validation error data to downstream functions.
	Results Results `yaml:"results,omitempty" json:"results,omitempty"`
}

// CheckResourceDuplication checks the GVKNN of resourceList.items to make sure they are unique. It returns errors if
// found more than one resource having the same GVKNN.
func CheckResourceDuplication(rl *ResourceList) error {
	idMap := map[yaml.ResourceIdentifier]struct{}{}
	for _, obj := range rl.Items {
		id := obj.resourceIdentifier()
		if _, ok := idMap[*id]; ok {
			return fmt.Errorf("duplicate Resource(apiVersion=%v, kind=%v, Namespace=%v, Name=%v)",
				obj.GetAPIVersion(), obj.GetKind(), obj.GetNamespace(), obj.GetName())
		}
		idMap[*id] = struct{}{}
	}
	return nil
}

// ParseResourceList parses a ResourceList from the input byte array. This function can be used to parse either KRM fn input
// or KRM fn output
func ParseResourceList(in []byte) (*ResourceList, error) {
	rl := &ResourceList{}
	rlObj, err := ParseKubeObject(in)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input bytes: %w", err)
	}
	if rlObj.GetKind() != kio.ResourceListKind {
		return nil, fmt.Errorf("input was of unexpected kind %q; expected ResourceList", rlObj.GetKind())
	}
	// Parse FunctionConfig. FunctionConfig can be empty, e.g. `kubeval` fn does not require a FunctionConfig.
	fc, found, err := rlObj.obj.GetNestedMap("functionConfig")
	if err != nil {
		return nil, fmt.Errorf("failed when tried to get functionConfig: %w", err)
	}
	if found {
		rl.FunctionConfig = asKubeObject(fc)
	} else {
		rl.FunctionConfig = NewEmptyKubeObject()
	}

	// Parse Items. Items can be empty, e.g. an input ResourceList for a generator function may not have items.
	items, found, err := rlObj.obj.GetNestedSlice("items")
	if err != nil {
		return nil, fmt.Errorf("failed when tried to get items: %w", err)
	}
	if found {
		objectItems, err := items.Elements()
		if err != nil {
			return nil, fmt.Errorf("failed to extract objects from items: %w", err)
		}
		for i := range objectItems {
			rl.Items = append(rl.Items, asKubeObject(objectItems[i]))
		}
	}

	// Parse Results. Results can be empty.
	res, found, err := rlObj.obj.GetNestedSlice("results")
	if err != nil {
		return nil, fmt.Errorf("failed when tried to get results: %w", err)
	}
	if found {
		var results Results
		err = res.Node().Decode(&results)
		if err != nil {
			return nil, fmt.Errorf("failed to decode results: %w", err)
		}
		rl.Results = results
	}
	return rl, nil
}

// toYNode converts the ResourceList to the yaml.Node representation.
func (rl *ResourceList) toYNode() (*yaml.Node, error) {
	reMap := internal.NewMap(nil)
	reObj := &KubeObject{SubObject{obj: reMap, parentGVK: schema.GroupVersionKind{}, fieldpath: ""}}
	if err := reObj.SetAPIVersion(kio.ResourceListAPIVersion); err != nil {
		return nil, err
	}
	if err := reObj.SetKind(kio.ResourceListKind); err != nil {
		return nil, err
	}

	if rl.Items != nil && len(rl.Items) > 0 {
		itemsSlice := internal.NewSliceVariant()
		for i := range rl.Items {
			itemsSlice.Add(rl.Items[i].node())
		}
		if err := reMap.SetNestedSlice(itemsSlice, "items"); err != nil {
			return nil, err
		}
	}
	if !rl.FunctionConfig.IsEmpty() {
		if err := reMap.SetNestedMap(rl.FunctionConfig.node(), "functionConfig"); err != nil {
			return nil, err
		}
	}

	if rl.Results != nil && len(rl.Results) > 0 {
		resultsSlice := internal.NewSliceVariant()
		for _, result := range rl.Results {
			mv, err := internal.TypedObjectToMapVariant(result)
			if err != nil {
				return nil, err
			}
			resultsSlice.Add(mv)
		}
		if err := reMap.SetNestedSlice(resultsSlice, "results"); err != nil {
			return nil, err
		}
	}

	return reMap.Node(), nil
}

// ToYAML converts the ResourceList to yaml.
func (rl *ResourceList) ToYAML() ([]byte, error) {
	// Sort the resources first.
	rl.Sort()
	ynode, err := rl.toYNode()
	if err != nil {
		return nil, err
	}
	doc := internal.NewDoc([]*yaml.Node{ynode}...)
	return doc.ToYAML()
}

// Sort sorts the ResourceList.items by apiVersion, kind, namespace and name.
func (rl *ResourceList) Sort() {
	sort.Sort(rl.Items)
}

// UpsertObjectToItems adds an object to ResourceList.items. The input object can
// be a KubeObject or any typed object (e.g. corev1.Pod).
func (rl *ResourceList) UpsertObjectToItems(obj interface{}, checkExistence func(obj, another *KubeObject) bool, replaceIfAlreadyExist bool) error {
	if checkExistence == nil {
		checkExistence = func(obj, another *KubeObject) bool {
			ri1 := obj.resourceIdentifier()
			ri2 := another.resourceIdentifier()
			return reflect.DeepEqual(ri1, ri2)
		}
	}

	var ko *KubeObject
	switch obj := obj.(type) {
	case KubeObject:
		ko = &obj
	case *KubeObject:
		ko = obj
	case yaml.RNode:
		ko = rnodeToKubeObject(&obj)
	case *yaml.RNode:
		ko = rnodeToKubeObject(obj)
	case yaml.Node:
		ko = rnodeToKubeObject(yaml.NewRNode(&obj))
	case *yaml.Node:
		ko = rnodeToKubeObject(yaml.NewRNode(obj))
	default:
		var err error
		ko, err = NewFromTypedObject(obj)
		if err != nil {
			return err
		}
	}

	idx := -1
	for i, item := range rl.Items {
		if checkExistence(ko, item) {
			idx = i
			break
		}
	}
	if idx == -1 {
		rl.Items = append(rl.Items, ko)
	} else if replaceIfAlreadyExist {
		rl.Items[idx] = ko
	}
	return nil
}

func (rl *ResourceList) LogResult(err error) {
	// If the error is not a Results type, we wrap the error as a Result.
	if err == nil {
		return
	}
	switch result := err.(type) {
	case Results:
		rl.Results = append(rl.Results, result...)
	case Result:
		rl.Results = append(rl.Results, &result)
	case *Result:
		rl.Results = append(rl.Results, result)
	default:
		rl.Results = append(rl.Results, ErrorResult(err))
	}
}

// ResourceListProcessor is implemented by configuration functions built with this framework
// to conform to the Configuration Functions Specification:
// https://github.com/kubernetes-sigs/kustomize/blob/master/cmd/config/docs/api-conventions/functions-spec.md
type ResourceListProcessor interface {
	Process(rl *ResourceList) (bool, error)
}

// ResourceListProcessorFunc converts a compatible function to a ResourceListProcessor.
type ResourceListProcessorFunc func(rl *ResourceList) (bool, error)

func (p ResourceListProcessorFunc) Process(rl *ResourceList) (bool, error) {
	return p(rl)
}

// Chain chains a list of ResourceListProcessor as a single ResourceListProcessor.
func Chain(processors ...ResourceListProcessor) ResourceListProcessor {
	return ResourceListProcessorFunc(func(rl *ResourceList) (bool, error) {
		success := true
		for _, processor := range processors {
			s, err := processor.Process(rl)
			if !s {
				success = false
			}
			if err != nil {
				return false, err
			}
		}
		return success, nil
	})
}

// ChainFunctions chains a list of ResourceListProcessorFunc as a single
// ResourceListProcessorFunc.
func ChainFunctions(functions ...ResourceListProcessorFunc) ResourceListProcessorFunc {
	return func(rl *ResourceList) (bool, error) {
		success := true
		for _, fn := range functions {
			s, err := fn(rl)
			if !s {
				success = false
			}
			if err != nil {
				return false, err
			}
		}
		return success, nil
	}
}

// ApplyFnBySelector iterates through every object in ResourceList.items, and if
// it satisfies the selector, fn will be applied on it.
func ApplyFnBySelector(rl *ResourceList, selector func(obj *KubeObject) bool, fn func(obj *KubeObject) error) error {
	var results Results
	for i, obj := range rl.Items {
		if !selector(obj) {
			continue
		}
		err := fn(rl.Items[i])
		if err == nil {
			continue
		}
		switch te := err.(type) {
		case Results:
			results = append(results, te...)
		case *Result:
			results = append(results, te)
		default:
			results = append(results, ErrorResult(err))
		}
	}
	if len(results) > 0 {
		rl.Results = results
		return results
	}
	return nil
}
