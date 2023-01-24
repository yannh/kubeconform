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

/*
Package fn provides the SDK to write KRM functions.

# Before you start

This fn SDK requires some basic KRM function Specification knowledge. To make the best usage of your time, we recommend
you to be familiar with "ResourceList" before moving forward.

	The KRM Function Specification, or "ResourceList", defines the standards of the inter-process communication between
	the orchestrator (i.e. kpt CLI) and functions.

See KRM Function Specification reference in https://github.com/kubernetes-sigs/kustomize/blob/master/cmd/config/docs/api-conventions/functions-spec.md

# KRM Function

A KRM function can mutate and/or validate Kubernetes resources in a ResourceList.

The ResourceList type and the KubeObject type are the core parts of this package.
The ResourceList type maps to the ResourceList in the function spec.

Read more about how to use KRM functions in https://kpt.dev/book/04-using-functions/

Read more about how to develop a KRM function in https://kpt.dev/book/05-developing-functions/

A general workflow is:
 1. Reads the "ResourceList" object from STDIN.
 2. Gets the function configs from the "ResourceList.FunctionConfig".
 3. Mutate or validate the Kubernetes YAML resources from the "ResourceList.Items" field with the function configs.
 4. Writes the modified "ResourceList" to STDOUT.
 5. Write function message to "ResourceList.Results" with severity "Info", "Warning" or "Error"

# KubeObject

The KubeObject is the basic unit to perform operations on KRM resources.

In the "AsMain", both "Items" and "FunctionConfig"
are converted to the KubeObject(s).

If you are familiar with unstructured.Unstructured, using KubeObject is as simple as using unstructured.Unstructured.
You can call function like `NestedStringOrDie` `SetNestedStringMap`, etc.

Except that KubeObject will not have pass-in interface arguments, nor will return an interface.
Instead, you shall treat each KubeObject field (slice, or non-string map）as SubObject.

SubObject also have most of the KubeObject methods, except the MetaType or NameType specific methods like "GetNamespace", "SetLabel".
This is because SubObject is designed as a sub object of KubeObject. SubObject to KubeObject is like `spec` section to `Deployment`.
You can get the Deployment name from `metadata.name`, KubeObject.GetName() or KubeObject.NestedString("metadata", "name").
But you cannot get "metadata.name" from a Deployment "spec". For "spec" SubObject, you can get the ".replicas" field by
SubObject.NestedInt64("replicas")

Besides unstructured style, another way to use KubeObject is to purely work on the KubeObject/SubObject by calling
"GetMap", "GetSlice", "UpsertMap" which expects the return to be SubObject(s) pointer.

# AsMain

"AsMain" is the main entrypoint. In most cases, you only need to provide the mutator or validation logic and have AsMain
handles the ResourceList parsing, KRM resource field type detection, read from STDIN and write to STDOUT.

"AsMain" accepts a struct that either implement the ResourceListProcessor interface or Runner interface.

See github.com/GoogleContainerTools/kpt-functions-sdk/go/fn/examples for detailed usage.
*/
package fn
