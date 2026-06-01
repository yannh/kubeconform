package sarif

// Edge - Represents a directed edge in a graph.
type Edge struct {
	// A string that uniquely identifies the edge within its graph.
	ID *string `json:"id,omitempty"`

	// A short description of the edge.
	Label *Message `json:"label,omitempty"`

	// Key/value pairs that provide additional information about the edge.
	Properties *PropertyBag `json:"properties,omitempty"`

	// Identifies the source node (the node at which the edge starts).
	SourceNodeID *string `json:"sourceNodeId,omitempty"`

	// Identifies the target node (the node at which the edge ends).
	TargetNodeID *string `json:"targetNodeId,omitempty"`
}

// NewEdge - creates a new
func NewEdge() *Edge {
	return &Edge{}
}

// WithID - add a ID to the Edge
func (i *Edge) WithID(id string) *Edge {
	i.ID = &id
	return i
}

// WithLabel - add a Label to the Edge
func (l *Edge) WithLabel(label *Message) *Edge {
	l.Label = label
	return l
}

// WithProperties - add a Properties to the Edge
func (p *Edge) WithProperties(properties *PropertyBag) *Edge {
	p.Properties = properties
	return p
}

// WithSourceNodeID - add a SourceNodeID to the Edge
func (s *Edge) WithSourceNodeID(sourceNodeId string) *Edge {
	s.SourceNodeID = &sourceNodeId
	return s
}

// WithTargetNodeID - add a TargetNodeID to the Edge
func (t *Edge) WithTargetNodeID(targetNodeId string) *Edge {
	t.TargetNodeID = &targetNodeId
	return t
}
