package sarif

import "github.com/google/uuid"

func NewGuid() Guid {
	return Guid(uuid.New().String())
}
