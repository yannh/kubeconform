package sarif

// GraphTraversal - Represents a path through a graph.
type GraphTraversal struct {
	// Values of relevant expressions at the start of the graph traversal that may change during graph traversal.
	InitialState map[string]MultiformatMessageString `json:"initialState,omitempty"`

	// Values of relevant expressions at the start of the graph traversal that remain constant for the graph traversal.
	ImmutableState map[string]MultiformatMessageString `json:"immutableState,omitempty"`

	// A description of this graph traversal.
	Description *Message `json:"description,omitempty"`

	// The sequences of edges traversed by this graph traversal.
	EdgeTraversals []*EdgeTraversal `json:"edgeTraversals"`

	// Key/value pairs that provide additional information about the graph traversal.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The index within the result.graphs to be associated with the result.
	ResultGraphIndex int `json:"resultGraphIndex"`

	// The index within the run.graphs to be associated with the result.
	RunGraphIndex int `json:"runGraphIndex"`
}

// NewGraphTraversal - creates a new
func NewGraphTraversal() *GraphTraversal {
	return &GraphTraversal{
		EdgeTraversals:   make([]*EdgeTraversal, 0),
		ResultGraphIndex: -1,
		RunGraphIndex:    -1,
	}
}

// AddInitialState - add a single InitialState to the GraphTraversal
func (i *GraphTraversal) AddInitialState(key string, initialState MultiformatMessageString) *GraphTraversal {
	i.InitialState[key] = initialState
	return i
}

// WithInitialState - add a InitialState to the GraphTraversal
func (i *GraphTraversal) WithInitialState(initialState map[string]MultiformatMessageString) *GraphTraversal {
	i.InitialState = initialState
	return i
}

// AddImmutableState - add a single ImmutableState to the GraphTraversal
func (i *GraphTraversal) AddImmutableState(key string, immutableState MultiformatMessageString) *GraphTraversal {
	i.ImmutableState[key] = immutableState
	return i
}

// WithImmutableState - add a ImmutableState to the GraphTraversal
func (i *GraphTraversal) WithImmutableState(immutableState map[string]MultiformatMessageString) *GraphTraversal {
	i.ImmutableState = immutableState
	return i
}

// WithDescription - add a Description to the GraphTraversal
func (d *GraphTraversal) WithDescription(description *Message) *GraphTraversal {
	d.Description = description
	return d
}

// WithEdgeTraversals - add a EdgeTraversals to the GraphTraversal
func (e *GraphTraversal) WithEdgeTraversals(edgeTraversals []*EdgeTraversal) *GraphTraversal {
	e.EdgeTraversals = edgeTraversals
	return e
}

// AddEdgeTraversal - add a single EdgeTraversal to the GraphTraversal
func (e *GraphTraversal) AddEdgeTraversal(edgeTraversal *EdgeTraversal) *GraphTraversal {
	e.EdgeTraversals = append(e.EdgeTraversals, edgeTraversal)
	return e
}

// WithProperties - add a Properties to the GraphTraversal
func (p *GraphTraversal) WithProperties(properties *PropertyBag) *GraphTraversal {
	p.Properties = properties
	return p
}

// WithResultGraphIndex - add a ResultGraphIndex to the GraphTraversal
func (r *GraphTraversal) WithResultGraphIndex(resultGraphIndex int) *GraphTraversal {
	r.ResultGraphIndex = resultGraphIndex
	return r
}

// WithRunGraphIndex - add a RunGraphIndex to the GraphTraversal
func (r *GraphTraversal) WithRunGraphIndex(runGraphIndex int) *GraphTraversal {
	r.RunGraphIndex = runGraphIndex
	return r
}
