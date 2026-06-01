package sarif

// ToolComponentReference - Identifies a particular toolComponent object, either the driver or an extension.
type ToolComponentReference struct {
	// The 'guid' property of the referenced toolComponent.
	GuID *string `json:"guid,omitempty"`

	// An index into the referenced toolComponent in tool.extensions.
	Index int `json:"index"`

	// The 'name' property of the referenced toolComponent.
	Name *string `json:"name,omitempty"`

	// Key/value pairs that provide additional information about the toolComponentReference.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewToolComponentReference - creates a new
func NewToolComponentReference() *ToolComponentReference {
	return &ToolComponentReference{
		Index: -1,
	}
}

// WithGuID - add a GuID to the ToolComponentReference
func (g *ToolComponentReference) WithGuID(guid string) *ToolComponentReference {
	g.GuID = &guid
	return g
}

// WithIndex - add a Index to the ToolComponentReference
func (i *ToolComponentReference) WithIndex(index int) *ToolComponentReference {
	i.Index = index
	return i
}

// WithName - add a Name to the ToolComponentReference
func (n *ToolComponentReference) WithName(name string) *ToolComponentReference {
	n.Name = &name
	return n
}

// WithProperties - add a Properties to the ToolComponentReference
func (p *ToolComponentReference) WithProperties(properties *PropertyBag) *ToolComponentReference {
	p.Properties = properties
	return p
}
