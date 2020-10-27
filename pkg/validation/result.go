package validation

import "github.com/yannh/kubeconform/pkg/resource"

// Status TODO
type Status int

// TODO
const (
	_ Status = iota
	Error
	Skipped
	Valid
	Invalid
	Empty
)

// Result TODO
type Result struct {
	resource.Resource
	Err    error
	Status Status
}
