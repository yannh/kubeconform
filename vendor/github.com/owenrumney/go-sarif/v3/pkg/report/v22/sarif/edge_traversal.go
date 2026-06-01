package sarif

// EdgeTraversal - Represents the traversal of a single edge during a graph traversal.
type EdgeTraversal struct {
	// The values of relevant expressions after the edge has been traversed.
	FinalState map[string]MultiformatMessageString `json:"finalState,omitempty"`

	// Identifies the edge being traversed.
	EdgeID *string `json:"edgeId,omitempty"`

	// A message to display to the user as the edge is traversed.
	Message *Message `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the edge traversal.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The number of edge traversals necessary to return from a nested graph.
	StepOverEdgeCount *int `json:"stepOverEdgeCount,omitempty"`
}

// NewEdgeTraversal - creates a new
func NewEdgeTraversal() *EdgeTraversal {
	return &EdgeTraversal{}
}

// AddFinalState - add a single FinalState to the EdgeTraversal
func (f *EdgeTraversal) AddFinalState(key string, finalState MultiformatMessageString) *EdgeTraversal {
	f.FinalState[key] = finalState
	return f
}

// WithFinalState - add a FinalState to the EdgeTraversal
func (f *EdgeTraversal) WithFinalState(finalState map[string]MultiformatMessageString) *EdgeTraversal {
	f.FinalState = finalState
	return f
}

// WithEdgeID - add a EdgeID to the EdgeTraversal
func (e *EdgeTraversal) WithEdgeID(edgeId string) *EdgeTraversal {
	e.EdgeID = &edgeId
	return e
}

// WithMessage - add a Message to the EdgeTraversal
func (m *EdgeTraversal) WithMessage(message *Message) *EdgeTraversal {
	m.Message = message
	return m
}

// WithProperties - add a Properties to the EdgeTraversal
func (p *EdgeTraversal) WithProperties(properties *PropertyBag) *EdgeTraversal {
	p.Properties = properties
	return p
}

// WithStepOverEdgeCount - add a StepOverEdgeCount to the EdgeTraversal
func (s *EdgeTraversal) WithStepOverEdgeCount(stepOverEdgeCount int) *EdgeTraversal {
	s.StepOverEdgeCount = &stepOverEdgeCount
	return s
}
