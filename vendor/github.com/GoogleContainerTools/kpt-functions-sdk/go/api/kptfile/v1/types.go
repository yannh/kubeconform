// Copyright 2021 Google LLC
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

// Package defines Kptfile schema.
// Version: v1
// swagger:meta
package v1

import (
	"fmt"
)

const (
	KptFileName       = "Kptfile"
	KptFileKind       = "Kptfile"
	KptFileGroup      = "kpt.dev"
	KptFileVersion    = "v1"
	KptFileAPIVersion = KptFileGroup + "/" + KptFileVersion
)

// ResourceMeta contains the metadata for a both Resource Type and Resource.
type ResourceMeta struct {
	TypeMeta `json:",inline" yaml:",inline"`
	// ObjectMeta is the metadata field of a Resource
	ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// TypeMeta partially copies apimachinery/pkg/apis/meta/v1.TypeMeta
// No need for a direct dependence; the fields are stable.
type TypeMeta struct {
	// APIVersion is the apiVersion field of a Resource
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	// Kind is the kind field of a Resource
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

// ObjectMeta contains metadata about a Resource
type ObjectMeta struct {
	NameMeta `json:",inline" yaml:",inline"`
	// Labels is the metadata.labels field of a Resource
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	// Annotations is the metadata.annotations field of a Resource.
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// NameMeta contains name information.
type NameMeta struct {
	// Name is the metadata.name field of a Resource
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Namespace is the metadata.namespace field of a Resource
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// KptFile contains information about a package managed with kpt.
// swagger:model kptfile
type KptFile struct {
	ResourceMeta `yaml:",inline" json:",inline"`

	Upstream *Upstream `yaml:"upstream,omitempty" json:"upstream,omitempty"`

	// UpstreamLock is a resolved locator for the last fetch of the package.
	UpstreamLock *UpstreamLock `yaml:"upstreamLock,omitempty" json:"upstreamLock,omitempty"`

	// Info contains metadata such as license, documentation, etc.
	Info *PackageInfo `yaml:"info,omitempty" json:"info,omitempty"`

	// Pipeline declares the pipeline of functions.
	Pipeline *Pipeline `yaml:"pipeline,omitempty" json:"pipeline,omitempty"`

	// Inventory contains parameters for the inventory object used in apply.
	Inventory *Inventory `yaml:"inventory,omitempty" json:"inventory,omitempty"`
}

// OriginType defines the type of origin for a package.
type OriginType string

const (
	// GitOrigin specifies a package as having been cloned from a git repository.
	GitOrigin OriginType = "git"
)

// UpdateStrategyType defines the strategy for updating a package from upstream.
type UpdateStrategyType string

// ToUpdateStrategy takes a string representing an update strategy and will
// return the strategy as an UpdateStrategyType. If the provided string does
// not match any known update strategies, an error will be returned.
func ToUpdateStrategy(strategy string) (UpdateStrategyType, error) {
	switch strategy {
	case string(ResourceMerge):
		return ResourceMerge, nil
	case string(FastForward):
		return FastForward, nil
	case string(ForceDeleteReplace):
		return ForceDeleteReplace, nil
	default:
		return "", fmt.Errorf("unknown update strategy %q", strategy)
	}
}

const (
	// ResourceMerge performs a structural schema-aware comparison and
	// merges the changes into the local package.
	ResourceMerge UpdateStrategyType = "resource-merge"
	// FastForward fails without updating if the local package was modified
	// since it was fetched.
	FastForward UpdateStrategyType = "fast-forward"
	// ForceDeleteReplace wipes all local changes to the package.
	ForceDeleteReplace UpdateStrategyType = "force-delete-replace"
)

// UpdateStrategies is a slice with all the supported update strategies.
var UpdateStrategies = []UpdateStrategyType{
	ResourceMerge,
	FastForward,
	ForceDeleteReplace,
}

// UpdateStrategiesAsStrings returns a list of update strategies as strings.
func UpdateStrategiesAsStrings() []string {
	var strs []string
	for _, s := range UpdateStrategies {
		strs = append(strs, string(s))
	}
	return strs
}

// Upstream is a user-specified upstream locator for a package.
type Upstream struct {
	// Type is the type of origin.
	Type OriginType `yaml:"type,omitempty" json:"type,omitempty"`

	// Git is the locator for a package stored on Git.
	Git *Git `yaml:"git,omitempty" json:"git,omitempty"`

	// UpdateStrategy declares how a package will be updated from upstream.
	UpdateStrategy UpdateStrategyType `yaml:"updateStrategy,omitempty" json:"updateStrategy,omitempty"`
}

// Git is the user-specified locator for a package on Git.
type Git struct {
	// Repo is the git repository the package.
	// e.g. 'https://github.com/kubernetes/examples.git'
	Repo string `yaml:"repo,omitempty" json:"repo,omitempty"`

	// Directory is the sub directory of the git repository.
	// e.g. 'staging/cockroachdb'
	Directory string `yaml:"directory,omitempty" json:"directory,omitempty"`

	// Ref can be a Git branch, tag, or a commit SHA-1.
	Ref string `yaml:"ref,omitempty" json:"ref,omitempty"`
}

// UpstreamLock is a resolved locator for the last fetch of the package.
type UpstreamLock struct {
	// Type is the type of origin.
	Type OriginType `yaml:"type,omitempty" json:"type,omitempty"`

	// Git is the resolved locator for a package on Git.
	Git *GitLock `yaml:"git,omitempty" json:"git,omitempty"`
}

// GitLock is the resolved locator for a package on Git.
type GitLock struct {
	// Repo is the git repository that was fetched.
	// e.g. 'https://github.com/kubernetes/examples.git'
	Repo string `yaml:"repo,omitempty" json:"repo,omitempty"`

	// Directory is the sub directory of the git repository that was fetched.
	// e.g. 'staging/cockroachdb'
	Directory string `yaml:"directory,omitempty" json:"directory,omitempty"`

	// Ref can be a Git branch, tag, or a commit SHA-1 that was fetched.
	// e.g. 'master'
	Ref string `yaml:"ref,omitempty" json:"ref,omitempty"`

	// Commit is the SHA-1 for the last fetch of the package.
	// This is set by kpt for bookkeeping purposes.
	Commit string `yaml:"commit,omitempty" json:"commit,omitempty"`
}

// PackageInfo contains optional information about the package such as license, documentation, etc.
// These fields are not consumed by any functionality in kpt and are simply passed through.
// Note that like any other KRM resource, humans and automation can also use `metadata.labels` and
// `metadata.annotations` as the extension mechanism.
type PackageInfo struct {
	// Site is the URL for package web page.
	Site string `yaml:"site,omitempty" json:"site,omitempty"`

	// Email is the list of emails for the package authors.
	Emails []string `yaml:"emails,omitempty" json:"emails,omitempty"`

	// SPDX license identifier (e.g. "Apache-2.0"). See: https://spdx.org/licenses/
	License string `yaml:"license,omitempty" json:"license,omitempty"`

	// Relative slash-delimited path to the license file (e.g. LICENSE.txt)
	LicenseFile string `yaml:"licenseFile,omitempty" json:"licenseFile,omitempty"`

	// Description contains a short description of the package.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`

	// Keywords is a list of keywords for this package.
	Keywords []string `yaml:"keywords,omitempty" json:"keywords,omitempty"`

	// Man is the path to documentation about the package
	Man string `yaml:"man,omitempty" json:"man,omitempty"`
}

// Subpackages declares a local or remote subpackage.
type Subpackage struct {
	// Name of the immediate subdirectory relative to this Kptfile where the suppackage
	// either exists (local subpackages) or will be fetched to (remote subpckages).
	// This must be unique across all subpckages of a package.
	LocalDir string `yaml:"localDir,omitempty" json:"localDir,omitempty"`

	// Upstream is a reference to where the subpackage should be fetched from.
	// Whether a subpackage is local or remote is determined by whether Upstream is specified.
	Upstream *Upstream `yaml:"upstream,omitempty" json:"upstream,omitempty"`
}

// Pipeline declares a pipeline of functions used to mutate or validate resources.
type Pipeline struct {
	//  Sources defines the source packages to resolve as input to the pipeline. Possible values:
	//  a) A slash-separated, OS-agnostic relative package path which may include '.' and '..' e.g. './base', '../foo'
	//     The source package is resolved recursively.
	//  b) Resources in this package using '.'. Meta resources such as the Kptfile, Pipeline, and function configs
	//     are excluded.
	//  c) Resources in this package AND all resolved subpackages using './*'
	//
	// Resultant list of resources are ordered:
	// - According to the order of sources specified in this array.
	// - When using './*': Subpackages are resolved in alphanumerical order before package resources.
	//
	// When omitted, defaults to './*'.
	// Sources []string `yaml:"sources,omitempty"`

	// Following fields define the sequence of functions in the pipeline.
	// Input of the first function is the resolved sources.
	// Input of the second function is the output of the first function, and so on.
	// Order of operation: mutators, validators

	// Mutators defines a list of of KRM functions that mutate resources.
	Mutators []Function `yaml:"mutators,omitempty" json:"mutators,omitempty"`

	// Validators defines a list of KRM functions that validate resources.
	// Validators are not permitted to mutate resources.
	Validators []Function `yaml:"validators,omitempty" json:"validators,omitempty"`
}

// String returns the string representation of Pipeline struct
// The string returned is the struct content in Go default format.
func (p *Pipeline) String() string {
	return fmt.Sprintf("%+v", *p)
}

// IsEmpty returns true if the pipeline doesn't contain any functions in any of
// the function chains (mutators, validators).
func (p *Pipeline) IsEmpty() bool {
	if p == nil {
		return true
	}
	if len(p.Mutators) == 0 && len(p.Validators) == 0 {
		return true
	}
	return false
}

// Function specifies a KRM function.
type Function struct {
	// `Image` specifies the function container image.
	// It can either be fully qualified, e.g.:
	//
	//	image: gcr.io/kpt-fn/set-labels
	//
	// Optionally, kpt can be configured to use a image
	// registry host-path that will be used to resolve the image path in case
	// the image path is missing (Defaults to gcr.io/kpt-fn).
	// e.g. The following resolves to gcr.io/kpt-fn/set-labels:
	//
	//	image: set-labels
	Image string `yaml:"image,omitempty" json:"image,omitempty"`

	// Exec specifies the function binary executable.
	// The executable can be fully qualified or it must exists in the $PATH e.g:
	//
	// 	 exec: set-namespace
	// 	 exec: /usr/local/bin/my-custom-fn
	Exec string `yaml:"exec,omitempty" json:"exec,omitempty"`

	// `ConfigPath` specifies a slash-delimited relative path to a file in the current directory
	// containing a KRM resource used as the function config. This resource is
	// excluded when resolving 'sources', and as a result cannot be operated on
	// by the pipeline.
	ConfigPath string `yaml:"configPath,omitempty" json:"configPath,omitempty"`

	// `ConfigMap` is a convenient way to specify a function config of kind ConfigMap.
	ConfigMap map[string]string `yaml:"configMap,omitempty" json:"configMap,omitempty"`

	// `Name` is used to uniquely identify the function declaration
	// this is primarily used for merging function declaration with upstream counterparts
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// `Selectors` are used to specify resources on which the function should be executed
	// if not specified, all resources are selected
	Selectors []Selector `yaml:"selectors,omitempty" json:"selectors,omitempty"`

	// `Exclude` are used to specify resources on which the function should NOT be executed.
	// If not specified, all resources selected by `Selectors` are selected.
	Exclusions []Selector `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

// Selector specifies the selection criteria
// please update IsEmpty method if more properties are added
type Selector struct {
	// APIVersion of the target resources
	APIVersion string `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	// Kind of the target resources
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty"`
	// Name of the target resources
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	// Namespace of the target resources
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	// Labels on the target resources
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	// Annotations on the target resources
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

// IsEmpty returns true of none of the selection criteria is specified
func (s Selector) IsEmpty() bool {
	return s.APIVersion == "" &&
		s.Namespace == "" &&
		s.Name == "" &&
		s.Kind == "" &&
		len(s.Labels) == 0 &&
		len(s.Annotations) == 0
}

// Inventory encapsulates the parameters for the inventory resource applied to a cluster.
// All of the the parameters are required if any are set.
type Inventory struct {
	// Namespace for the inventory resource.
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	// Name of the inventory resource.
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	// Unique label to identify inventory resource in cluster.
	InventoryID string            `yaml:"inventoryID,omitempty" json:"inventoryID,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

func (i Inventory) IsValid() bool {
	// Name and Namespace are required inventory fields, so we check these 2 fields.
	// InventoryID is an optional field since we only store it locally if the user
	// specifies one.
	return i.Name != "" && i.Namespace != ""
}
