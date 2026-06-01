package sarif

// ReportingDescriptorRelationship - Information about the relation of one reporting descriptor to another.
type ReportingDescriptorRelationship struct {
	// A description of the reporting descriptor relationship.
	Description *Message `json:"description,omitempty"`

	// A set of distinct strings that categorize the relationship. Well-known kinds include 'canPrecede', 'canFollow', 'willPrecede', 'willFollow', 'superset', 'subset', 'equal', 'disjoint', 'relevant', and 'incomparable'.
	Kinds []string `json:"kinds"`

	// Key/value pairs that provide additional information about the reporting descriptor reference.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A reference to the related reporting descriptor.
	Target *ReportingDescriptorReference `json:"target,omitempty"`
}

// NewReportingDescriptorRelationship - creates a new
func NewReportingDescriptorRelationship() *ReportingDescriptorRelationship {
	return &ReportingDescriptorRelationship{
		Kinds: []string{"relevant"},
	}
}

// WithDescription - add a Description to the ReportingDescriptorRelationship
func (d *ReportingDescriptorRelationship) WithDescription(description *Message) *ReportingDescriptorRelationship {
	d.Description = description
	return d
}

// WithKinds - add a Kinds to the ReportingDescriptorRelationship
func (k *ReportingDescriptorRelationship) WithKinds(kinds []string) *ReportingDescriptorRelationship {
	k.Kinds = kinds
	return k
}

// AddKind - add a single Kind to the ReportingDescriptorRelationship
func (k *ReportingDescriptorRelationship) AddKind(kind string) *ReportingDescriptorRelationship {
	k.Kinds = append(k.Kinds, kind)
	return k
}

// WithProperties - add a Properties to the ReportingDescriptorRelationship
func (p *ReportingDescriptorRelationship) WithProperties(properties *PropertyBag) *ReportingDescriptorRelationship {
	p.Properties = properties
	return p
}

// WithTarget - add a Target to the ReportingDescriptorRelationship
func (t *ReportingDescriptorRelationship) WithTarget(target *ReportingDescriptorReference) *ReportingDescriptorRelationship {
	t.Target = target
	return t
}
