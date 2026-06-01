package sarif

// LogicalLocation - A logical location of a construct that produced a result.
type LogicalLocation struct {
	// The machine-readable name for the logical location, such as a mangled function name provided by a C++ compiler that encodes calling convention, return type and other details along with the function name.
	DecoratedName *string `json:"decoratedName,omitempty"`

	// The human-readable fully qualified name of the logical location.
	FullyQualifiedName *string `json:"fullyQualifiedName,omitempty"`

	// The index within the logical locations array.
	Index int `json:"index"`

	// The type of construct this logical location component refers to. Should be one of 'function', 'member', 'module', 'namespace', 'parameter', 'resource', 'returnType', 'type', 'variable', 'object', 'array', 'property', 'value', 'element', 'text', 'attribute', 'comment', 'declaration', 'dtd' or 'processingInstruction', if any of those accurately describe the construct.
	Kind *string `json:"kind,omitempty"`

	// Identifies the construct in which the result occurred. For example, this property might contain the name of a class or a method.
	Name *string `json:"name,omitempty"`

	// Identifies the index of the immediate parent of the construct in which the result was detected. For example, this property might point to a logical location that represents the namespace that holds a type.
	ParentIndex int `json:"parentIndex"`

	// Key/value pairs that provide additional information about the logical location.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewLogicalLocation - creates a new
func NewLogicalLocation() *LogicalLocation {
	return &LogicalLocation{
		Index:       -1,
		ParentIndex: -1,
	}
}

// WithDecoratedName - add a DecoratedName to the LogicalLocation
func (d *LogicalLocation) WithDecoratedName(decoratedName string) *LogicalLocation {
	d.DecoratedName = &decoratedName
	return d
}

// WithFullyQualifiedName - add a FullyQualifiedName to the LogicalLocation
func (f *LogicalLocation) WithFullyQualifiedName(fullyQualifiedName string) *LogicalLocation {
	f.FullyQualifiedName = &fullyQualifiedName
	return f
}

// WithIndex - add a Index to the LogicalLocation
func (i *LogicalLocation) WithIndex(index int) *LogicalLocation {
	i.Index = index
	return i
}

// WithKind - add a Kind to the LogicalLocation
func (k *LogicalLocation) WithKind(kind string) *LogicalLocation {
	k.Kind = &kind
	return k
}

// WithName - add a Name to the LogicalLocation
func (n *LogicalLocation) WithName(name string) *LogicalLocation {
	n.Name = &name
	return n
}

// WithParentIndex - add a ParentIndex to the LogicalLocation
func (p *LogicalLocation) WithParentIndex(parentIndex int) *LogicalLocation {
	p.ParentIndex = parentIndex
	return p
}

// WithProperties - add a Properties to the LogicalLocation
func (p *LogicalLocation) WithProperties(properties *PropertyBag) *LogicalLocation {
	p.Properties = properties
	return p
}
