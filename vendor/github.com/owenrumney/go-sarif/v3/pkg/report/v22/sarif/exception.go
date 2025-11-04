package sarif

// Exception - Describes a runtime exception encountered during the execution of an analysis tool.
type Exception struct {
	// An array of exception objects each of which is considered a cause of this exception.
	InnerExceptions []*Exception `json:"innerExceptions"`

	// A string that identifies the kind of exception, for example, the fully qualified type name of an object that was thrown, or the symbolic name of a signal.
	Kind *string `json:"kind,omitempty"`

	// A message that describes the exception.
	Message *string `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the exception.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The sequence of function calls leading to the exception.
	Stack *Stack `json:"stack,omitempty"`
}

// NewException - creates a new
func NewException() *Exception {
	return &Exception{
		InnerExceptions: make([]*Exception, 0),
	}
}

// WithInnerExceptions - add a InnerExceptions to the Exception
func (i *Exception) WithInnerExceptions(innerExceptions []*Exception) *Exception {
	i.InnerExceptions = innerExceptions
	return i
}

// AddInnerException - add a single InnerException to the Exception
func (i *Exception) AddInnerException(innerException *Exception) *Exception {
	i.InnerExceptions = append(i.InnerExceptions, innerException)
	return i
}

// WithKind - add a Kind to the Exception
func (k *Exception) WithKind(kind string) *Exception {
	k.Kind = &kind
	return k
}

// WithMessage - add a Message to the Exception
func (m *Exception) WithMessage(message string) *Exception {
	m.Message = &message
	return m
}

// WithProperties - add a Properties to the Exception
func (p *Exception) WithProperties(properties *PropertyBag) *Exception {
	p.Properties = properties
	return p
}

// WithStack - add a Stack to the Exception
func (s *Exception) WithStack(stack *Stack) *Exception {
	s.Stack = stack
	return s
}
