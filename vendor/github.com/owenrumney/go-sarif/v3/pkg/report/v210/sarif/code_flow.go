package sarif

// CodeFlow - A set of threadFlows which together describe a pattern of code execution relevant to detecting a result.
type CodeFlow struct {
	// A message relevant to the code flow.
	Message *Message `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the code flow.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of one or more unique threadFlow objects, each of which describes the progress of a program through a thread of execution.
	ThreadFlows []*ThreadFlow `json:"threadFlows,omitempty"`
}

// NewCodeFlow - creates a new
func NewCodeFlow() *CodeFlow {
	return &CodeFlow{
		ThreadFlows: make([]*ThreadFlow, 0),
	}
}

// WithMessage - add a Message to the CodeFlow
func (m *CodeFlow) WithMessage(message *Message) *CodeFlow {
	m.Message = message
	return m
}

// WithProperties - add a Properties to the CodeFlow
func (p *CodeFlow) WithProperties(properties *PropertyBag) *CodeFlow {
	p.Properties = properties
	return p
}

// WithThreadFlows - add a ThreadFlows to the CodeFlow
func (t *CodeFlow) WithThreadFlows(threadFlows []*ThreadFlow) *CodeFlow {
	t.ThreadFlows = threadFlows
	return t
}

// AddThreadFlow - add a single ThreadFlow to the CodeFlow
func (t *CodeFlow) AddThreadFlow(threadFlow *ThreadFlow) *CodeFlow {
	t.ThreadFlows = append(t.ThreadFlows, threadFlow)
	return t
}
