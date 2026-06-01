package sarif

// Attachment - An artifact relevant to a result.
type Attachment struct {
	// The location of the attachment.
	ArtifactLocation *ArtifactLocation `json:"artifactLocation,omitempty"`

	// A message describing the role played by the attachment.
	Description *Message `json:"description,omitempty"`

	// Key/value pairs that provide additional information about the attachment.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of rectangles specifying areas of interest within the image.
	Rectangles []*Rectangle `json:"rectangles"`

	// An array of regions of interest within the attachment.
	Regions []*Region `json:"regions"`
}

// NewAttachment - creates a new
func NewAttachment() *Attachment {
	return &Attachment{
		Rectangles: make([]*Rectangle, 0),
		Regions:    make([]*Region, 0),
	}
}

// WithArtifactLocation - add a ArtifactLocation to the Attachment
func (a *Attachment) WithArtifactLocation(artifactLocation *ArtifactLocation) *Attachment {
	a.ArtifactLocation = artifactLocation
	return a
}

// WithDescription - add a Description to the Attachment
func (d *Attachment) WithDescription(description *Message) *Attachment {
	d.Description = description
	return d
}

// WithProperties - add a Properties to the Attachment
func (p *Attachment) WithProperties(properties *PropertyBag) *Attachment {
	p.Properties = properties
	return p
}

// WithRectangles - add a Rectangles to the Attachment
func (r *Attachment) WithRectangles(rectangles []*Rectangle) *Attachment {
	r.Rectangles = rectangles
	return r
}

// AddRectangle - add a single Rectangle to the Attachment
func (r *Attachment) AddRectangle(rectangle *Rectangle) *Attachment {
	r.Rectangles = append(r.Rectangles, rectangle)
	return r
}

// WithRegions - add a Regions to the Attachment
func (r *Attachment) WithRegions(regions []*Region) *Attachment {
	r.Regions = regions
	return r
}

// AddRegion - add a single Region to the Attachment
func (r *Attachment) AddRegion(region *Region) *Attachment {
	r.Regions = append(r.Regions, region)
	return r
}
