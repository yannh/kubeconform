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
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func WithContext(ctx context.Context, runner Runner) ResourceListProcessor {
	return runnerProcessor{ctx: ctx, fnRunner: runner}
}

type runnerProcessor struct {
	ctx      context.Context
	fnRunner Runner
}

// EmptyFunctionConfig is a workaround solution to handle the case where kpt passes in a functionConfig placeholder
// (Configmap with empty `data`) if user does not provide the actual FunctionConfig. Ideally, kpt should pass in an empty
// FunctionConfig object.
func EmptyFunctionConfig(o *KubeObject) bool {
	data, found, err := o.NestedStringMap("data")
	// Some other type
	if !found || err != nil {
		return false
	}
	return o.GetKind() == "ConfigMap" && o.GetName() == "function-input" && len(data) == 0
}

// Process assigns the ResourceList.FunctionConfig to Runner's attributes, and calls the Runner.Run methods (main method)
// to run functions. The r.fnRunner accepts three kinds of functionConfig value:
//  1. no function config, it only runs fnRunner.Run
//  2. ConfigMap type, it requires the Runner instance to have one contributes of type map[string]string to receive the ConfigMap `.data` value.
//  3. Runner type, it uses the Runner struct name as the FunctionConfig Kind. e.g. if the Runner is `SetNamespace`,
//     the FunctionConfig should be `{"Kind": "SetNamespace", "apiVersion": "fn.kpt.dev/v1alpha1"}
func (r runnerProcessor) Process(rl *ResourceList) (bool, error) {
	// Validate and Parse the input FunctionConfig to r.fnRunner
	if rl.FunctionConfig.IsEmpty() || EmptyFunctionConfig(rl.FunctionConfig) {
		// functions may not need functionConfig.
		rl.Results.Infof("`FunctionConfig` is not given")
	} else {
		err := r.config(rl.FunctionConfig)
		if err != nil {
			rl.Results.ErrorE(err)
			return false, nil
		}
	}
	// Run the main function.
	fnCtx := &Context{Context: r.ctx}
	results := new(Results)
	shouldPass := r.fnRunner.Run(fnCtx, rl.FunctionConfig, rl.Items, results)
	// If running in a pipeline, the ResourceList may already have results from previous function runs.
	// Thus, we only append new results to the end.
	rl.Results = append(rl.Results, *results...)
	return shouldPass, nil
}

func (r *runnerProcessor) config(o *KubeObject) error {
	if r.fnRunner == nil {
		return fmt.Errorf("the object which implements `ResourceListProcessor` interface requires a `Runner` or `fnRunner` attribute," +
			" got nil")
	}
	switch o.GroupKind() {
	case schema.GroupKind{Kind: "ConfigMap"}:
		data, _, err := o.NestedStringMap("data")
		if data == nil {
			return err
		}
		return assignCMDataToFn(r.fnRunner, data)
	case schema.GroupKind{Group: KptFunctionGroup, Kind: asFnName(r.fnRunner)}:
		return o.As(r.fnRunner)
	default:
		return fmt.Errorf("unknown FunctionConfig `%v`, expect `%v.%v` or `ConfigMap.v1`", o.GroupKind(), asFnName(r.fnRunner), KptFunctionGroup)
	}
}

func asFnName(runner Runner) string {
	// Validate the fnRunner type to avoid panic.
	kind := reflect.ValueOf(runner).Kind()
	if kind != reflect.Interface && kind != reflect.Ptr {
		return ""
	}
	return reflect.ValueOf(runner).Elem().Type().Name()
}

func assignCMDataToFn(runner Runner, data map[string]string) error {
	obj := reflect.ValueOf(runner).Elem()
	if obj.Kind() != reflect.Struct {
		return fmt.Errorf("the ConfigMap is not of a struct, got %v", obj.Kind().String())
	}
	stringMap := reflect.MapOf(reflect.TypeOf("string"), reflect.TypeOf("string"))
	for i := 0; i < obj.NumField(); i++ {
		if obj.Field(i).Kind() == reflect.Map && obj.Field(i).Type() == stringMap {
			if obj.Field(i).CanSet() {
				obj.Field(i).Set(reflect.ValueOf(data))
			}
			return nil
		}
	}
	return fmt.Errorf("unable to assign the given ConfigMap `.data` to FunctionConfig %v. please make sure the %v "+
		"has a field of type map[string]string", asFnName(runner), asFnName(runner))
}
