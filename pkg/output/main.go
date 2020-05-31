package output

import (
	"github.com/yannh/kubeconform/pkg/validator"
)

const (
	VALID   = iota
	INVALID = iota
	ERROR   = iota
	SKIPPED = iota
)

type Output interface {
	Write(filename string, err error, skipped bool)
	Flush()
}

func status(err error, skipped bool) int {
	if skipped {
		return SKIPPED
	}

	if err != nil {
		if _, ok := err.(validator.InvalidResourceError); ok {
			return INVALID
		} else {
			return ERROR
		}
	}

	return VALID
}
