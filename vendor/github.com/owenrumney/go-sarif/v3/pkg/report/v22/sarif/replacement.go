package sarif

// Replacement - The replacement of a single region of an artifact.
type Replacement struct {
	// The region of the artifact to delete.
	DeletedRegion *Region `json:"deletedRegion,omitempty"`

	// The content to insert at the location specified by the 'deletedRegion' property.
	InsertedContent *ArtifactContent `json:"insertedContent,omitempty"`

	// Key/value pairs that provide additional information about the replacement.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewReplacement - creates a new
func NewReplacement() *Replacement {
	return &Replacement{}
}

// WithDeletedRegion - add a DeletedRegion to the Replacement
func (d *Replacement) WithDeletedRegion(deletedRegion *Region) *Replacement {
	d.DeletedRegion = deletedRegion
	return d
}

// WithInsertedContent - add a InsertedContent to the Replacement
func (i *Replacement) WithInsertedContent(insertedContent *ArtifactContent) *Replacement {
	i.InsertedContent = insertedContent
	return i
}

// WithProperties - add a Properties to the Replacement
func (p *Replacement) WithProperties(properties *PropertyBag) *Replacement {
	p.Properties = properties
	return p
}
