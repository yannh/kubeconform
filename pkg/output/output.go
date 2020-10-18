package output

import (
	"github.com/yannh/kubeconform/pkg/validator"
)

const (
	_ = iota
	VALID
	INVALID
	ERROR
	SKIPPED
	EMPTY
)

type Output interface {
	Write(filename, kind, name, version string, err error, skipped bool) error
	Flush() error
}

func status(kind, name string, err error, skipped bool) int {
	if name == "" && kind == "" && err == nil && skipped == false {
		return EMPTY
	}

	if skipped {
		return SKIPPED
	}

	if err != nil {
		if _, ok := err.(validator.InvalidResourceError); ok {
			return INVALID
		}
		return ERROR
	}

	return VALID
}
