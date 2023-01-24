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

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	// upstreamIdentifierRegexPattern provides the rough regex to parse a upstream-identiifier annotation.
	// "group" should be a domain name. We accept empty string for kubernetes core v1 resources.
	// "kind" should be the resource type with initial in capitals.
	// "namespace" should follow RFC 1123 Label Names. We accept "~C~ for cluster-scoped resource or unknown scope resources.
	// "name" should follow RFC 1123 Label Names https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
	upstreamIdentifierRegexPattern = `(?P<group>[a-z0-9-.]*)\|(?P<kind>[A-Z][a-zA-Z0-9]*)\|(?P<namespace>[a-z0-9-]{1,63}|~C)\|(?P<name>[a-z0-9-]{1,63})`
	upstreamIdentifierFormat       = "<GROUP>|<KIND>|<NAMESPACE>|<NAME>"
	regexPatternGroup              = "group"
	regexPatternKind               = "kind"
	regexPatterNamespace           = "namespace"
	regexPatternName               = "name"
)

type ResourceIdentifier struct {
	Group     string
	Version   string
	Kind      string
	Name      string
	Namespace string
}

func (r *ResourceIdentifier) String() string {
	return fmt.Sprintf("%v|%v|%v|%v", r.Group, r.Kind, r.Namespace, r.Name)
}

// hasUpstreamIdentifier determines whether the args are touching the kpt only annotation "internal.kpt.dev/upstream-identifier"
func (o *SubObject) hasUpstreamIdentifier(val interface{}, fields ...string) bool {
	kind := reflect.ValueOf(val).Kind()
	if kind == reflect.Ptr {
		kind = reflect.TypeOf(val).Elem().Kind()
	}
	switch kind {
	case reflect.String:
		if fields[len(fields)-1] == UpstreamIdentifier {
			return true
		}
	case reflect.Map:
		if fields[len(fields)-1] == "annotations" {
			for _, key := range reflect.ValueOf(val).MapKeys() {
				if key.String() == UpstreamIdentifier {
					return true
				}
			}
		}
	}
	return false
}

func (o *KubeObject) effectiveNamespace() string {
	if o.HasNamespace() {
		return o.GetNamespace()
	}
	if o.IsNamespaceScoped() {
		return DefaultNamespace
	}
	return UnknownNamespace
}

// GetId gets the Group, Kind, Namespace and Name as the ResourceIdentifier.
func (o *KubeObject) GetId() *ResourceIdentifier {
	group, _ := ParseGroupVersion(o.GetAPIVersion())
	return &ResourceIdentifier{
		Group:     group,
		Kind:      o.GetKind(),
		Namespace: o.effectiveNamespace(),
		Name:      o.GetName(),
	}
}

func parseUpstreamIdentifier(upstreamId string) (*ResourceIdentifier, error) {
	upstreamId = strings.TrimSpace(upstreamId)
	r := regexp.MustCompile(upstreamIdentifierRegexPattern)
	match := r.FindStringSubmatch(upstreamId)
	if match == nil {
		return nil, &ErrInternalAnnotation{Message: fmt.Sprintf("annotation %v: %v is in bad format. expect %q",
			UpstreamIdentifier, upstreamId, upstreamIdentifierFormat)}
	}
	matchGroups := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			matchGroups[name] = match[i]
		}
	}
	return &ResourceIdentifier{
		Group:     matchGroups[regexPatternGroup],
		Kind:      matchGroups[regexPatternKind],
		Namespace: matchGroups[regexPatterNamespace],
		Name:      matchGroups[regexPatternName],
	}, nil
}

// GetOriginId provides the `ResourceIdentifier` to identify the upstream origin of a KRM resource.
// This origin is generated and maintained by kpt pkg management and is stored in the `internal.kpt.dev/upstream-identiifer` annotation.
// If a resource does not have an upstream origin, we use its current meta resource ID instead.
func (o *KubeObject) GetOriginId() (*ResourceIdentifier, error) {
	upstreamId := o.GetAnnotation(UpstreamIdentifier)
	if upstreamId != "" {
		return parseUpstreamIdentifier(upstreamId)
	}
	return o.GetId(), nil
}

// HasUpstreamOrigin tells whether a resource is sourced from an upstream package resource.
func (o *KubeObject) HasUpstreamOrigin() bool {
	upstreamId := o.GetAnnotation(UpstreamIdentifier)
	return upstreamId != ""
}

// ParseGroupVersion parses a "apiVersion" to get the "group" and "version" values.
func ParseGroupVersion(apiVersion string) (group, version string) {
	if i := strings.Index(apiVersion, "/"); i > -1 {
		return apiVersion[:i], apiVersion[i+1:]
	}
	return "", apiVersion
}
