package sarif

// Address - A physical or virtual address, or a range of addresses, in an 'addressable region' (memory or a binary file).
type Address struct {
	// The address expressed as a byte offset from the start of the addressable region.
	AbsoluteAddress int `json:"absoluteAddress"`

	// A human-readable fully qualified name that is associated with the address.
	FullyQualifiedName *string `json:"fullyQualifiedName,omitempty"`

	// The index within run.addresses of the cached object for this address.
	Index int `json:"index"`

	// An open-ended string that identifies the address kind. 'data', 'function', 'header','instruction', 'module', 'page', 'section', 'segment', 'stack', 'stackFrame', 'table' are well-known values.
	Kind *string `json:"kind,omitempty"`

	// The number of bytes in this range of addresses.
	Length *int `json:"length,omitempty"`

	// A name that is associated with the address, e.g., '.text'.
	Name *string `json:"name,omitempty"`

	// The byte offset of this address from the absolute or relative address of the parent object.
	OffsetFromParent *int `json:"offsetFromParent,omitempty"`

	// The index within run.addresses of the parent object.
	ParentIndex int `json:"parentIndex"`

	// Key/value pairs that provide additional information about the address.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The address expressed as a byte offset from the absolute address of the top-most parent object.
	RelativeAddress *int `json:"relativeAddress,omitempty"`
}

// NewAddress - creates a new
func NewAddress() *Address {
	return &Address{
		AbsoluteAddress: -1,
		Index:           -1,
		ParentIndex:     -1,
	}
}

// WithAbsoluteAddress - add a AbsoluteAddress to the Address
func (a *Address) WithAbsoluteAddress(absoluteAddress int) *Address {
	a.AbsoluteAddress = absoluteAddress
	return a
}

// WithFullyQualifiedName - add a FullyQualifiedName to the Address
func (f *Address) WithFullyQualifiedName(fullyQualifiedName string) *Address {
	f.FullyQualifiedName = &fullyQualifiedName
	return f
}

// WithIndex - add a Index to the Address
func (i *Address) WithIndex(index int) *Address {
	i.Index = index
	return i
}

// WithKind - add a Kind to the Address
func (k *Address) WithKind(kind string) *Address {
	k.Kind = &kind
	return k
}

// WithLength - add a Length to the Address
func (l *Address) WithLength(length int) *Address {
	l.Length = &length
	return l
}

// WithName - add a Name to the Address
func (n *Address) WithName(name string) *Address {
	n.Name = &name
	return n
}

// WithOffsetFromParent - add a OffsetFromParent to the Address
func (o *Address) WithOffsetFromParent(offsetFromParent int) *Address {
	o.OffsetFromParent = &offsetFromParent
	return o
}

// WithParentIndex - add a ParentIndex to the Address
func (p *Address) WithParentIndex(parentIndex int) *Address {
	p.ParentIndex = parentIndex
	return p
}

// WithProperties - add a Properties to the Address
func (p *Address) WithProperties(properties *PropertyBag) *Address {
	p.Properties = properties
	return p
}

// WithRelativeAddress - add a RelativeAddress to the Address
func (r *Address) WithRelativeAddress(relativeAddress int) *Address {
	r.RelativeAddress = &relativeAddress
	return r
}
