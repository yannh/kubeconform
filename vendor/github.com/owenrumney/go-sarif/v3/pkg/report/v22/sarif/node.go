package sarif

// Node - Represents a node in a graph.
type Node struct {
	// Array of child nodes.
	Children []*Node `json:"children"`

	// A string that uniquely identifies the node within its graph.
	ID *string `json:"id,omitempty"`

	// A short description of the node.
	Label *Message `json:"label,omitempty"`

	// A code location associated with the node.
	Location *Location `json:"location,omitempty"`

	// Key/value pairs that provide additional information about the node.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewNode - creates a new
func NewNode() *Node {
	return &Node{
		Children: make([]*Node, 0),
	}
}

// WithChildren - add a Children to the Node
func (c *Node) WithChildren(children []*Node) *Node {
	c.Children = children
	return c
}

// AddChildren - add a single Children to the Node
func (c *Node) AddChildren(children *Node) *Node {
	c.Children = append(c.Children, children)
	return c
}

// WithID - add a ID to the Node
func (i *Node) WithID(id string) *Node {
	i.ID = &id
	return i
}

// WithLabel - add a Label to the Node
func (l *Node) WithLabel(label *Message) *Node {
	l.Label = label
	return l
}

// WithLocation - add a Location to the Node
func (l *Node) WithLocation(location *Location) *Node {
	l.Location = location
	return l
}

// WithProperties - add a Properties to the Node
func (p *Node) WithProperties(properties *PropertyBag) *Node {
	p.Properties = properties
	return p
}
