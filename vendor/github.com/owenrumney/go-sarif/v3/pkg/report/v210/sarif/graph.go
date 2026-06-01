package sarif

// Graph - A network of nodes and directed edges that describes some aspect of the structure of the code (for example, a call graph).
type Graph struct {
	// A description of the graph.
	Description *Message `json:"description,omitempty"`

	// An array of edge objects representing the edges of the graph.
	Edges []*Edge `json:"edges"`

	// An array of node objects representing the nodes of the graph.
	Nodes []*Node `json:"nodes"`

	// Key/value pairs that provide additional information about the graph.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewGraph - creates a new
func NewGraph() *Graph {
	return &Graph{
		Edges: make([]*Edge, 0),
		Nodes: make([]*Node, 0),
	}
}

// WithDescription - add a Description to the Graph
func (d *Graph) WithDescription(description *Message) *Graph {
	d.Description = description
	return d
}

// WithEdges - add a Edges to the Graph
func (e *Graph) WithEdges(edges []*Edge) *Graph {
	e.Edges = edges
	return e
}

// AddEdge - add a single Edge to the Graph
func (e *Graph) AddEdge(edge *Edge) *Graph {
	e.Edges = append(e.Edges, edge)
	return e
}

// WithNodes - add a Nodes to the Graph
func (n *Graph) WithNodes(nodes []*Node) *Graph {
	n.Nodes = nodes
	return n
}

// AddNode - add a single Node to the Graph
func (n *Graph) AddNode(node *Node) *Graph {
	n.Nodes = append(n.Nodes, node)
	return n
}

// WithProperties - add a Properties to the Graph
func (p *Graph) WithProperties(properties *PropertyBag) *Graph {
	p.Properties = properties
	return p
}
