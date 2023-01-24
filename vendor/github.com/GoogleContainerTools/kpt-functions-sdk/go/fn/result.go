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
	"fmt"
	"sort"
	"strings"
)

// Severity indicates the severity of the Result
type Severity string

const (
	// Error indicates the result is an error.  Will cause the function to exit non-0.
	Error Severity = "error"
	// Warning indicates the result is a warning
	Warning Severity = "warning"
	// Info indicates the result is an informative message
	Info Severity = "info"
)

// ResourceRef fills the ResourceRef field in Results
type ResourceRef struct {
	APIVersion string `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	Kind       string `yaml:"kind,omitempty" json:"kind,omitempty"`
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
	Namespace  string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}

// Result defines a validation result
type Result struct {
	// Message is a human readable message. This field is required.
	Message string `yaml:"message,omitempty" json:"message,omitempty"`

	// Severity is the severity of this result
	Severity Severity `yaml:"severity,omitempty" json:"severity,omitempty"`

	// ResourceRef is a reference to a resource.
	// Required fields: apiVersion, kind, name.
	ResourceRef *ResourceRef `yaml:"resourceRef,omitempty" json:"resourceRef,omitempty"`

	// Field is a reference to the field in a resource this result refers to
	Field *Field `yaml:"field,omitempty" json:"field,omitempty"`

	// File references a file containing the resource this result refers to
	File *File `yaml:"file,omitempty" json:"file,omitempty"`

	// Tags is an unstructured key value map stored with a result that may be set
	// by external tools to store and retrieve arbitrary metadata
	Tags map[string]string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

func (i Result) Error() string {
	return (i).String()
}

// String provides a human-readable message for the result item
func (i Result) String() string {
	identifier := i.ResourceRef
	var idStringList []string
	if identifier != nil {
		if identifier.APIVersion != "" {
			idStringList = append(idStringList, identifier.APIVersion)
		}
		if identifier.Kind != "" {
			idStringList = append(idStringList, identifier.Kind)
		}
		if identifier.Namespace != "" {
			idStringList = append(idStringList, identifier.Namespace)
		}
		if identifier.Name != "" {
			idStringList = append(idStringList, identifier.Name)
		}
	}
	formatString := "[%s]"
	severity := i.Severity
	// We default Severity to Info when converting a result to a message.
	if i.Severity == "" {
		severity = Info
	}
	list := []interface{}{severity}
	if len(idStringList) > 0 {
		formatString += " %s"
		list = append(list, strings.Join(idStringList, "/"))
	}
	if i.Field != nil {
		formatString += " %s"
		list = append(list, i.Field.Path)
	}
	formatString += ": %s"
	list = append(list, i.Message)
	return fmt.Sprintf(formatString, list...)
}

// File references a file containing a resource
type File struct {
	// Path is relative path to the file containing the resource.
	// This field is required.
	Path string `yaml:"path,omitempty" json:"path,omitempty"`

	// Index is the index into the file containing the resource
	// (i.e. if there are multiple resources in a single file)
	Index int `yaml:"index,omitempty" json:"index,omitempty"`
}

// Field references a field in a resource
type Field struct {
	// Path is the field path. This field is required.
	Path string `yaml:"path,omitempty" json:"path,omitempty"`

	// CurrentValue is the current field value
	CurrentValue interface{} `yaml:"currentValue,omitempty" json:"currentValue,omitempty"`

	// ProposedValue is the proposed value of the field to fix an issue.
	ProposedValue interface{} `yaml:"proposedValue,omitempty" json:"proposedValue,omitempty"`
}

type Results []*Result

// Errorf writes an Error level `result` to the results slice. It accepts arguments according to a format specifier.
// e.g.
// results.Errorf("bad kind %v", "invalid")
func (r *Results) Errorf(format string, a ...any) {
	errResult := &Result{Severity: Error, Message: fmt.Sprintf(format, a...)}
	*r = append(*r, errResult)
}

// ErrorE writes the `error` as an Error level `result` to the results slice.
// e.g.
//
//	err := error.New("test)
//	results.ErrorE(err)
func (r *Results) ErrorE(err error) {
	errResult := &Result{Severity: Error, Message: err.Error()}
	*r = append(*r, errResult)
}

// Infof writes an Info level `result` to the results slice. It accepts arguments according to a format specifier.
// e.g.
//
//	results.Infof("update %v %q ", "ConfigMap", "kptfile.kpt.dev")
func (r *Results) Infof(format string, a ...any) {
	infoResult := &Result{Severity: Info, Message: fmt.Sprintf(format, a...)}
	*r = append(*r, infoResult)
}

// Warningf writes a Warning level `result` to the results slice. It accepts arguments according to a format specifier.
// e.g.
//
//	results.Warningf("bad kind %q", "invalid")
func (r *Results) Warningf(format string, a ...any) {
	warnResult := &Result{Severity: Warning, Message: fmt.Sprintf(format, a...)}
	*r = append(*r, warnResult)
}

// WarningE writes an error as a Warning level `result` to the results slice.
// Normally this function can be used for cases that need error tolerance.
func (r *Results) WarningE(err error) {
	warnResult := &Result{Severity: Warning, Message: err.Error()}
	*r = append(*r, warnResult)
}

func (r *Results) String() string {
	var results []string
	for _, result := range *r {
		results = append(results, result.String())
	}
	return strings.Join(results, "\n---\n")
}

// Error enables Results to be returned as an error
func (r Results) Error() string {
	var msgs []string
	for _, i := range r {
		msgs = append(msgs, i.String())
	}
	return strings.Join(msgs, "\n\n")
}

// ExitCode provides the exit code based on the result's severity
func (r Results) ExitCode() int {
	for _, i := range r {
		if i.Severity == Error {
			return 1
		}
	}
	return 0
}

// Sort performs an in place stable sort of Results
func (r Results) Sort() {
	sort.SliceStable(r, func(i, j int) bool {
		if fileLess(r, i, j) != 0 {
			return fileLess(r, i, j) < 0
		}
		if severityLess(r, i, j) != 0 {
			return severityLess(r, i, j) < 0
		}
		return resultToString(*r[i]) < resultToString(*r[j])
	})
}

func severityLess(items Results, i, j int) int {
	severityToNumber := map[Severity]int{
		Error:   0,
		Warning: 1,
		Info:    2,
	}

	severityLevelI, found := severityToNumber[items[i].Severity]
	if !found {
		severityLevelI = 3
	}
	severityLevelJ, found := severityToNumber[items[j].Severity]
	if !found {
		severityLevelJ = 3
	}
	return severityLevelI - severityLevelJ
}

func fileLess(items Results, i, j int) int {
	var fileI, fileJ File
	if items[i].File == nil {
		fileI = File{}
	} else {
		fileI = *items[i].File
	}
	if items[j].File == nil {
		fileJ = File{}
	} else {
		fileJ = *items[j].File
	}
	if fileI.Path != fileJ.Path {
		if fileI.Path < fileJ.Path {
			return -1
		}
		return 1
	}
	return fileI.Index - fileJ.Index
}

func resultToString(item Result) string {
	return fmt.Sprintf("resource-ref:%s,field:%s,message:%s",
		item.ResourceRef, item.Field, item.Message)
}

func ErrorConfigFileResult(err error, path string) *Result {
	return ConfigFileResult(err.Error(), path, Error)
}

func ConfigFileResult(msg, path string, severity Severity) *Result {
	return &Result{
		Message:  msg,
		Severity: severity,
		File: &File{
			Path: path,
		},
	}
}

func ErrorResult(err error) *Result {
	return GeneralResult(err.Error(), Error)
}

func GeneralResult(msg string, severity Severity) *Result {
	return &Result{
		Message:  msg,
		Severity: severity,
	}
}

func ErrorConfigObjectResult(err error, obj *KubeObject) *Result {
	return ConfigObjectResult(err.Error(), obj, Error)
}

func ConfigObjectResult(msg string, obj *KubeObject, severity Severity) *Result {
	return &Result{
		Message:  msg,
		Severity: severity,
		ResourceRef: &ResourceRef{
			APIVersion: obj.GetAPIVersion(),
			Kind:       obj.GetKind(),
			Name:       obj.GetName(),
			Namespace:  obj.GetNamespace(),
		},
		File: &File{
			Path:  obj.PathAnnotation(),
			Index: obj.IndexAnnotation(),
		},
	}
}
