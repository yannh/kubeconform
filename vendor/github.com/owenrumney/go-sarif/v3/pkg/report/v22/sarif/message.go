package sarif

// Message - Encapsulates a message intended to be read by the end user.
type Message struct {
	// An array of strings to substitute into the message string.
	Arguments []string `json:"arguments"`

	// The identifier for this message.
	ID *string `json:"id,omitempty"`

	// A Markdown message string.
	Markdown *string `json:"markdown,omitempty"`

	// Key/value pairs that provide additional information about the message.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A plain text message string.
	Text *string `json:"text,omitempty"`
}

// NewMessage - creates a new
func NewMessage() *Message {
	return &Message{
		Arguments: make([]string, 0),
	}
}

// WithArguments - add a Arguments to the Message
func (a *Message) WithArguments(arguments []string) *Message {
	a.Arguments = arguments
	return a
}

// AddArgument - add a single Argument to the Message
func (a *Message) AddArgument(argument string) *Message {
	a.Arguments = append(a.Arguments, argument)
	return a
}

// WithID - add a ID to the Message
func (i *Message) WithID(id string) *Message {
	i.ID = &id
	return i
}

// WithMarkdown - add a Markdown to the Message
func (m *Message) WithMarkdown(markdown string) *Message {
	m.Markdown = &markdown
	return m
}

// WithProperties - add a Properties to the Message
func (p *Message) WithProperties(properties *PropertyBag) *Message {
	p.Properties = properties
	return p
}

// WithText - add a Text to the Message
func (t *Message) WithText(text string) *Message {
	t.Text = &text
	return t
}
