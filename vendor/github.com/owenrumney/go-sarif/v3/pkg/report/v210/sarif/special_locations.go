package sarif

// SpecialLocations - Defines locations of special significance to SARIF consumers.
type SpecialLocations struct {
	// Provides a suggestion to SARIF consumers to display file paths relative to the specified location.
	DisplayBase *ArtifactLocation `json:"displayBase,omitempty"`

	// Key/value pairs that provide additional information about the special locations.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewSpecialLocations - creates a new
func NewSpecialLocations() *SpecialLocations {
	return &SpecialLocations{}
}

// WithDisplayBase - add a DisplayBase to the SpecialLocations
func (d *SpecialLocations) WithDisplayBase(displayBase *ArtifactLocation) *SpecialLocations {
	d.DisplayBase = displayBase
	return d
}

// WithProperties - add a Properties to the SpecialLocations
func (p *SpecialLocations) WithProperties(properties *PropertyBag) *SpecialLocations {
	p.Properties = properties
	return p
}
