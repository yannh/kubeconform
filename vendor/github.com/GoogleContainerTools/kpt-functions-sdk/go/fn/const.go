// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package fn

const (
	// internalPrefix is the prefix given to internal annotations that are used
	// internally by the orchestrator
	internalPrefix string = "internal.config.kubernetes.io/"

	// IndexAnnotation records the index of a specific resource in a file or input stream.
	IndexAnnotation string = internalPrefix + "index"

	// PathAnnotation records the path to the file the Resource was read from
	PathAnnotation string = internalPrefix + "path"

	// SeqIndentAnnotation records the sequence nodes indentation of the input resource
	SeqIndentAnnotation string = internalPrefix + "seqindent"

	// IdAnnotation records the id of the resource to map inputs to outputs
	IdAnnotation string = internalPrefix + "id"

	// InternalAnnotationsMigrationResourceIDAnnotation is used to uniquely identify
	// resources during round trip to and from a function execution. We will use it
	// to track the internal annotations and reconcile them if needed.
	InternalAnnotationsMigrationResourceIDAnnotation = internalPrefix + "annotations-migration-resource-id"

	// ConfigPrefix is the prefix given to the custom kubernetes annotations.
	ConfigPrefix string = "config.kubernetes.io/"

	// KptLocalConfig marks a KRM resource to be skipped from deploying to the cluster via `kpt live apply`.
	KptLocalConfig = ConfigPrefix + "local-config"
)

// For Kpt use only constants
const (
	// KptUseOnlyPrefix is the prefix of kpt-only annotations. Users are not expected to touch these annotations.
	KptUseOnlyPrefix = "internal.kpt.dev/"

	// UpstreamIdentifier is the annotation to record a resource's upstream origin.
	// It is in the form of <GROUP>|<KIND>|<NAMESPACE>|<NAME>
	UpstreamIdentifier = KptUseOnlyPrefix + "upstream-identifier"

	// UnknownNamespace is the special char for cluster-scoped or unknown-scoped resources. This is only used in upstream-identifier
	UnknownNamespace = "~C"
	// DefaultNamespace is the actual namespace value if a namespace-scoped resource has its namespace field unspecified.
	DefaultNamespace = "default"
)

// For KPT Function Configuration
const (
	// KptFunctionGroup is the group name for the KRM resource which defines the configuration of a function execution.
	// See KRM function specification `ResourceList.FunctionConfig`
	KptFunctionGroup = "fn.kpt.dev"
	// KptFunctionGroup is the version for the KRM resource which defines the configuration of a function execution.
	// See KRM function specification `ResourceList.FunctionConfig`
	KptFunctionVersion = "v1alpha1"
	// KptFunctionGroup is the ApiVersion for the KRM resource which defines the configuration of a function execution.
	// See KRM function specification `ResourceList.FunctionConfig`
	KptFunctionApiVersion = KptFunctionGroup + "/" + KptFunctionVersion
)
