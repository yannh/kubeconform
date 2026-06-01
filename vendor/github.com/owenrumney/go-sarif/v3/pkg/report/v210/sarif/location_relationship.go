package sarif

// LocationRelationship - Information about the relation of one location to another.
type LocationRelationship struct {
	// A description of the location relationship.
	Description *Message `json:"description,omitempty"`

	// A set of distinct strings that categorize the relationship. Well-known kinds include 'includes', 'isIncludedBy' and 'relevant'.
	Kinds []string `json:"kinds"`

	// Key/value pairs that provide additional information about the location relationship.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A reference to the related location.
	Target *int `json:"target,omitempty"`
}

// NewLocationRelationship - creates a new
func NewLocationRelationship() *LocationRelationship {
	return &LocationRelationship{
		Kinds: []string{"relevant"},
	}
}

// WithDescription - add a Description to the LocationRelationship
func (d *LocationRelationship) WithDescription(description *Message) *LocationRelationship {
	d.Description = description
	return d
}

// WithKinds - add a Kinds to the LocationRelationship
func (k *LocationRelationship) WithKinds(kinds []string) *LocationRelationship {
	k.Kinds = kinds
	return k
}

// AddKind - add a single Kind to the LocationRelationship
func (k *LocationRelationship) AddKind(kind string) *LocationRelationship {
	k.Kinds = append(k.Kinds, kind)
	return k
}

// WithProperties - add a Properties to the LocationRelationship
func (p *LocationRelationship) WithProperties(properties *PropertyBag) *LocationRelationship {
	p.Properties = properties
	return p
}

// WithTarget - add a Target to the LocationRelationship
func (t *LocationRelationship) WithTarget(target int) *LocationRelationship {
	t.Target = &target
	return t
}
