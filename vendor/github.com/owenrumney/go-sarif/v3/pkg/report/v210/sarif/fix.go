package sarif

// Fix - A proposed fix for the problem represented by a result object. A fix specifies a set of artifacts to modify. For each artifact, it specifies a set of bytes to remove, and provides a set of new bytes to replace them.
type Fix struct {
	// One or more artifact changes that comprise a fix for a result.
	ArtifactChanges []*ArtifactChange `json:"artifactChanges,omitempty"`

	// A message that describes the proposed fix, enabling viewers to present the proposed change to an end user.
	Description *Message `json:"description,omitempty"`

	// Key/value pairs that provide additional information about the fix.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewFix - creates a new
func NewFix() *Fix {
	return &Fix{
		ArtifactChanges: make([]*ArtifactChange, 0),
	}
}

// WithArtifactChanges - add a ArtifactChanges to the Fix
func (a *Fix) WithArtifactChanges(artifactChanges []*ArtifactChange) *Fix {
	a.ArtifactChanges = artifactChanges
	return a
}

// AddArtifactChange - add a single ArtifactChange to the Fix
func (a *Fix) AddArtifactChange(artifactChange *ArtifactChange) *Fix {
	a.ArtifactChanges = append(a.ArtifactChanges, artifactChange)
	return a
}

// WithDescription - add a Description to the Fix
func (d *Fix) WithDescription(description *Message) *Fix {
	d.Description = description
	return d
}

// WithProperties - add a Properties to the Fix
func (p *Fix) WithProperties(properties *PropertyBag) *Fix {
	p.Properties = properties
	return p
}
