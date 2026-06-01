package sarif

// Region - A region within an artifact where a result was detected.
type Region struct {
	// The length of the region in bytes.
	ByteLength *int `json:"byteLength,omitempty"`

	// The zero-based offset from the beginning of the artifact of the first byte in the region.
	ByteOffset int `json:"byteOffset"`

	// The length of the region in characters.
	CharLength *int `json:"charLength,omitempty"`

	// The zero-based offset from the beginning of the artifact of the first character in the region.
	CharOffset int `json:"charOffset"`

	// The column number of the character following the end of the region.
	EndColumn *int `json:"endColumn,omitempty"`

	// The line number of the last character in the region.
	EndLine *int `json:"endLine,omitempty"`

	// A message relevant to the region.
	Message *Message `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the region.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The portion of the artifact contents within the specified region.
	Snippet *ArtifactContent `json:"snippet,omitempty"`

	// Specifies the source language, if any, of the portion of the artifact specified by the region object.
	SourceLanguage *string `json:"sourceLanguage,omitempty"`

	// The column number of the first character in the region.
	StartColumn *int `json:"startColumn,omitempty"`

	// The line number of the first character in the region.
	StartLine *int `json:"startLine,omitempty"`
}

// NewRegion - creates a new
func NewRegion() *Region {
	return &Region{
		ByteOffset: -1,
		CharOffset: -1,
	}
}

// WithByteLength - add a ByteLength to the Region
func (b *Region) WithByteLength(byteLength int) *Region {
	b.ByteLength = &byteLength
	return b
}

// WithByteOffset - add a ByteOffset to the Region
func (b *Region) WithByteOffset(byteOffset int) *Region {
	b.ByteOffset = byteOffset
	return b
}

// WithCharLength - add a CharLength to the Region
func (c *Region) WithCharLength(charLength int) *Region {
	c.CharLength = &charLength
	return c
}

// WithCharOffset - add a CharOffset to the Region
func (c *Region) WithCharOffset(charOffset int) *Region {
	c.CharOffset = charOffset
	return c
}

// WithEndColumn - add a EndColumn to the Region
func (e *Region) WithEndColumn(endColumn int) *Region {
	e.EndColumn = &endColumn
	return e
}

// WithEndLine - add a EndLine to the Region
func (e *Region) WithEndLine(endLine int) *Region {
	e.EndLine = &endLine
	return e
}

// WithMessage - add a Message to the Region
func (m *Region) WithMessage(message *Message) *Region {
	m.Message = message
	return m
}

// WithProperties - add a Properties to the Region
func (p *Region) WithProperties(properties *PropertyBag) *Region {
	p.Properties = properties
	return p
}

// WithSnippet - add a Snippet to the Region
func (s *Region) WithSnippet(snippet *ArtifactContent) *Region {
	s.Snippet = snippet
	return s
}

// WithSourceLanguage - add a SourceLanguage to the Region
func (s *Region) WithSourceLanguage(sourceLanguage string) *Region {
	s.SourceLanguage = &sourceLanguage
	return s
}

// WithStartColumn - add a StartColumn to the Region
func (s *Region) WithStartColumn(startColumn int) *Region {
	s.StartColumn = &startColumn
	return s
}

// WithStartLine - add a StartLine to the Region
func (s *Region) WithStartLine(startLine int) *Region {
	s.StartLine = &startLine
	return s
}
