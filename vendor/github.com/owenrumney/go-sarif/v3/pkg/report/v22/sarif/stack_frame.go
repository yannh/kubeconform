package sarif

// StackFrame - A function call within a stack trace.
type StackFrame struct {
	// The location to which this stack frame refers.
	Location *Location `json:"location,omitempty"`

	// The name of the module that contains the code of this stack frame.
	Module *string `json:"module,omitempty"`

	// The parameters of the call that is executing.
	Parameters []string `json:"parameters"`

	// Key/value pairs that provide additional information about the stack frame.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The thread identifier of the stack frame.
	ThreadID *int `json:"threadId,omitempty"`
}

// NewStackFrame - creates a new
func NewStackFrame() *StackFrame {
	return &StackFrame{
		Parameters: make([]string, 0),
	}
}

// WithLocation - add a Location to the StackFrame
func (l *StackFrame) WithLocation(location *Location) *StackFrame {
	l.Location = location
	return l
}

// WithModule - add a Module to the StackFrame
func (m *StackFrame) WithModule(module string) *StackFrame {
	m.Module = &module
	return m
}

// WithParameters - add a Parameters to the StackFrame
func (p *StackFrame) WithParameters(parameters []string) *StackFrame {
	p.Parameters = parameters
	return p
}

// AddParameter - add a single Parameter to the StackFrame
func (p *StackFrame) AddParameter(parameter string) *StackFrame {
	p.Parameters = append(p.Parameters, parameter)
	return p
}

// WithProperties - add a Properties to the StackFrame
func (p *StackFrame) WithProperties(properties *PropertyBag) *StackFrame {
	p.Properties = properties
	return p
}

// WithThreadID - add a ThreadID to the StackFrame
func (t *StackFrame) WithThreadID(threadId int) *StackFrame {
	t.ThreadID = &threadId
	return t
}
