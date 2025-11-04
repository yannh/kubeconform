package sarif

// NewTextMessage creates a simple text message
func NewTextMessage(text string) *Message {
	return NewMessage().WithText(text)
}
