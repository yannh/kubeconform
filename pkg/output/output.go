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
)

type Output interface {
	Write(filename, kind, name, version string, err error, skipped bool) error
	Flush() error
}

func status(err error, skipped bool) int {
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
