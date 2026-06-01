package sarif

// NewSimpleRegion creates a new Region with the start and end line
func NewSimpleRegion(startLine, endLine int) *Region {
	return NewRegion().
		WithStartLine(startLine).
		WithEndLine(endLine)
}
