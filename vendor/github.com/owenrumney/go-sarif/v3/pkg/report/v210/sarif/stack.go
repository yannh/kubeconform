package sarif

// Stack - A call stack that is relevant to a result.
type Stack struct {
	// An array of stack frames that represents a sequence of calls, rendered in reverse chronological order, that comprise the call stack.
	Frames []*StackFrame `json:"frames,omitempty"`

	// A message relevant to this call stack.
	Message *Message `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the stack.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewStack - creates a new
func NewStack() *Stack {
	return &Stack{
		Frames: make([]*StackFrame, 0),
	}
}

// WithFrames - add a Frames to the Stack
func (f *Stack) WithFrames(frames []*StackFrame) *Stack {
	f.Frames = frames
	return f
}

// AddFrame - add a single Frame to the Stack
func (f *Stack) AddFrame(frame *StackFrame) *Stack {
	f.Frames = append(f.Frames, frame)
	return f
}

// WithMessage - add a Message to the Stack
func (m *Stack) WithMessage(message *Message) *Stack {
	m.Message = message
	return m
}

// WithProperties - add a Properties to the Stack
func (p *Stack) WithProperties(properties *PropertyBag) *Stack {
	p.Properties = properties
	return p
}
