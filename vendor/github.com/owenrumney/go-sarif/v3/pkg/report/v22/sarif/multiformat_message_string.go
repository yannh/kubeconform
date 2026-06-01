package sarif

// MultiformatMessageString - A message string or message format string rendered in multiple formats.
type MultiformatMessageString struct {
	// A Markdown message string or format string.
	Markdown *string `json:"markdown,omitempty"`

	// Key/value pairs that provide additional information about the message.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A plain text message string or format string.
	Text *string `json:"text,omitempty"`
}

// NewMultiformatMessageString - creates a new
func NewMultiformatMessageString() *MultiformatMessageString {
	return &MultiformatMessageString{}
}

// WithMarkdown - add a Markdown to the MultiformatMessageString
func (m *MultiformatMessageString) WithMarkdown(markdown string) *MultiformatMessageString {
	m.Markdown = &markdown
	return m
}

// WithProperties - add a Properties to the MultiformatMessageString
func (p *MultiformatMessageString) WithProperties(properties *PropertyBag) *MultiformatMessageString {
	p.Properties = properties
	return p
}

// WithText - add a Text to the MultiformatMessageString
func (t *MultiformatMessageString) WithText(text string) *MultiformatMessageString {
	t.Text = &text
	return t
}
