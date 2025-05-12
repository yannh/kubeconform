package loader

// NotFoundError is returned when the registry does not contain a schema for the resource
type NotFoundError struct {
	err error
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{err}
}
func (e *NotFoundError) Error() string   { return e.err.Error() }
func (e *NotFoundError) Retryable() bool { return false }

type NonJSONResponseError struct {
	err error
}

func NewNonJSONResponseError(err error) *NotFoundError {
	return &NotFoundError{err}
}
func (e *NonJSONResponseError) Error() string   { return e.err.Error() }
func (e *NonJSONResponseError) Retryable() bool { return false }
