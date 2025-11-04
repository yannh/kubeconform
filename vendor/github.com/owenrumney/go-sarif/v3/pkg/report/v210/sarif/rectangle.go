package sarif

// Rectangle - An area within an image.
type Rectangle struct {
	// The Y coordinate of the bottom edge of the rectangle, measured in the image's natural units.
	Bottom *float64 `json:"bottom,omitempty"`

	// The X coordinate of the left edge of the rectangle, measured in the image's natural units.
	Left *float64 `json:"left,omitempty"`

	// A message relevant to the rectangle.
	Message *Message `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the rectangle.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The X coordinate of the right edge of the rectangle, measured in the image's natural units.
	Right *float64 `json:"right,omitempty"`

	// The Y coordinate of the top edge of the rectangle, measured in the image's natural units.
	Top *float64 `json:"top,omitempty"`
}

// NewRectangle - creates a new
func NewRectangle() *Rectangle {
	return &Rectangle{}
}

// WithBottom - add a Bottom to the Rectangle
func (b *Rectangle) WithBottom(bottom float64) *Rectangle {
	b.Bottom = &bottom
	return b
}

// WithLeft - add a Left to the Rectangle
func (l *Rectangle) WithLeft(left float64) *Rectangle {
	l.Left = &left
	return l
}

// WithMessage - add a Message to the Rectangle
func (m *Rectangle) WithMessage(message *Message) *Rectangle {
	m.Message = message
	return m
}

// WithProperties - add a Properties to the Rectangle
func (p *Rectangle) WithProperties(properties *PropertyBag) *Rectangle {
	p.Properties = properties
	return p
}

// WithRight - add a Right to the Rectangle
func (r *Rectangle) WithRight(right float64) *Rectangle {
	r.Right = &right
	return r
}

// WithTop - add a Top to the Rectangle
func (t *Rectangle) WithTop(top float64) *Rectangle {
	t.Top = &top
	return t
}
