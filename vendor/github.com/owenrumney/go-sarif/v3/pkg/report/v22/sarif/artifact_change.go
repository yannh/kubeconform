package sarif

// ArtifactChange - A change to a single artifact.
type ArtifactChange struct {
	// The location of the artifact to change.
	ArtifactLocation *ArtifactLocation `json:"artifactLocation,omitempty"`

	// Key/value pairs that provide additional information about the change.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of replacement objects, each of which represents the replacement of a single region in a single artifact specified by 'artifactLocation'.
	Replacements []*Replacement `json:"replacements,omitempty"`
}

// NewArtifactChange - creates a new
func NewArtifactChange() *ArtifactChange {
	return &ArtifactChange{
		Replacements: make([]*Replacement, 0),
	}
}

// WithArtifactLocation - add a ArtifactLocation to the ArtifactChange
func (a *ArtifactChange) WithArtifactLocation(artifactLocation *ArtifactLocation) *ArtifactChange {
	a.ArtifactLocation = artifactLocation
	return a
}

// WithProperties - add a Properties to the ArtifactChange
func (p *ArtifactChange) WithProperties(properties *PropertyBag) *ArtifactChange {
	p.Properties = properties
	return p
}

// WithReplacements - add a Replacements to the ArtifactChange
func (r *ArtifactChange) WithReplacements(replacements []*Replacement) *ArtifactChange {
	r.Replacements = replacements
	return r
}

// AddReplacement - add a single Replacement to the ArtifactChange
func (r *ArtifactChange) AddReplacement(replacement *Replacement) *ArtifactChange {
	r.Replacements = append(r.Replacements, replacement)
	return r
}
