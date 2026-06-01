package sarif

// PhysicalLocation - A physical location relevant to a result. Specifies a reference to a programming artifact together with a range of bytes or characters within that artifact.
type PhysicalLocation struct {
	// The address of the location.
	Address *Address `json:"address,omitempty"`

	// The location of the artifact.
	ArtifactLocation *ArtifactLocation `json:"artifactLocation,omitempty"`

	// Specifies a portion of the artifact that encloses the region. Allows a viewer to display additional context around the region.
	ContextRegion *Region `json:"contextRegion,omitempty"`

	// Key/value pairs that provide additional information about the physical location.
	Properties *PropertyBag `json:"properties,omitempty"`

	// Specifies a portion of the artifact.
	Region *Region `json:"region,omitempty"`
}

// NewPhysicalLocation - creates a new
func NewPhysicalLocation() *PhysicalLocation {
	return &PhysicalLocation{}
}

// WithAddress - add a Address to the PhysicalLocation
func (a *PhysicalLocation) WithAddress(address *Address) *PhysicalLocation {
	a.Address = address
	return a
}

// WithArtifactLocation - add a ArtifactLocation to the PhysicalLocation
func (a *PhysicalLocation) WithArtifactLocation(artifactLocation *ArtifactLocation) *PhysicalLocation {
	a.ArtifactLocation = artifactLocation
	return a
}

// WithContextRegion - add a ContextRegion to the PhysicalLocation
func (c *PhysicalLocation) WithContextRegion(contextRegion *Region) *PhysicalLocation {
	c.ContextRegion = contextRegion
	return c
}

// WithProperties - add a Properties to the PhysicalLocation
func (p *PhysicalLocation) WithProperties(properties *PropertyBag) *PhysicalLocation {
	p.Properties = properties
	return p
}

// WithRegion - add a Region to the PhysicalLocation
func (r *PhysicalLocation) WithRegion(region *Region) *PhysicalLocation {
	r.Region = region
	return r
}
