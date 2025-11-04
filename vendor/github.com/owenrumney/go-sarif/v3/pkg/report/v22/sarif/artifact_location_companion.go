package sarif

// NewSimpleArtifactLocation creates a new SimpleArtifactLocation and returns a pointer to it
func NewSimpleArtifactLocation(uri string) *ArtifactLocation {
	return NewArtifactLocation().WithURI(uri)
}
