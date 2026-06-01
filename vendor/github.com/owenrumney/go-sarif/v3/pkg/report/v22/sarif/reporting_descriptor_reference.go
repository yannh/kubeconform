package sarif

// ReportingDescriptorReference - Information about how to locate a relevant reporting descriptor.
type ReportingDescriptorReference struct {
	// A guid that uniquely identifies the descriptor.
	Guid *Guid `json:"guid,omitempty"`

	// The id of the descriptor.
	ID *string `json:"id,omitempty"`

	// The index into an array of descriptors in toolComponent.ruleDescriptors, toolComponent.notificationDescriptors, or toolComponent.taxonomyDescriptors, depending on context.
	Index int `json:"index"`

	// Key/value pairs that provide additional information about the reporting descriptor reference.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A reference used to locate the toolComponent associated with the descriptor.
	ToolComponent *ToolComponentReference `json:"toolComponent,omitempty"`
}

// NewReportingDescriptorReference - creates a new
func NewReportingDescriptorReference() *ReportingDescriptorReference {
	return &ReportingDescriptorReference{
		Index: -1,
	}
}

// WithGuid - add a Guid to the ReportingDescriptorReference
func (g *ReportingDescriptorReference) WithGuid(guid *Guid) *ReportingDescriptorReference {
	g.Guid = guid
	return g
}

// WithID - add a ID to the ReportingDescriptorReference
func (i *ReportingDescriptorReference) WithID(id string) *ReportingDescriptorReference {
	i.ID = &id
	return i
}

// WithIndex - add a Index to the ReportingDescriptorReference
func (i *ReportingDescriptorReference) WithIndex(index int) *ReportingDescriptorReference {
	i.Index = index
	return i
}

// WithProperties - add a Properties to the ReportingDescriptorReference
func (p *ReportingDescriptorReference) WithProperties(properties *PropertyBag) *ReportingDescriptorReference {
	p.Properties = properties
	return p
}

// WithToolComponent - add a ToolComponent to the ReportingDescriptorReference
func (t *ReportingDescriptorReference) WithToolComponent(toolComponent *ToolComponentReference) *ReportingDescriptorReference {
	t.ToolComponent = toolComponent
	return t
}
