package sarif

// ArtifactLocation - Specifies the location of an artifact.
type ArtifactLocation struct {
	// A short description of the artifact location.
	Description *Message `json:"description,omitempty"`

	// The index within the run artifacts array of the artifact object associated with the artifact location.
	Index int `json:"index"`

	// Key/value pairs that provide additional information about the artifact location.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A string containing a valid relative or absolute URI.
	URI *string `json:"uri,omitempty"`

	// A string which indirectly specifies the absolute URI with respect to which a relative URI in the "uri" property is interpreted.
	URIBaseID *string `json:"uriBaseId,omitempty"`
}

// NewArtifactLocation - creates a new
func NewArtifactLocation() *ArtifactLocation {
	return &ArtifactLocation{
		Index: -1,
	}
}

// WithDescription - add a Description to the ArtifactLocation
func (d *ArtifactLocation) WithDescription(description *Message) *ArtifactLocation {
	d.Description = description
	return d
}

// WithIndex - add a Index to the ArtifactLocation
func (i *ArtifactLocation) WithIndex(index int) *ArtifactLocation {
	i.Index = index
	return i
}

// WithProperties - add a Properties to the ArtifactLocation
func (p *ArtifactLocation) WithProperties(properties *PropertyBag) *ArtifactLocation {
	p.Properties = properties
	return p
}

// WithURI - add a URI to the ArtifactLocation
func (u *ArtifactLocation) WithURI(uri string) *ArtifactLocation {
	u.URI = &uri
	return u
}

// WithURIBaseID - add a URIBaseID to the ArtifactLocation
func (u *ArtifactLocation) WithURIBaseID(uriBaseId string) *ArtifactLocation {
	u.URIBaseID = &uriBaseId
	return u
}
