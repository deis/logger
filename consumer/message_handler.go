package consumer

// MessageHandler is the interface to handle a single message
type MessageHandler interface {
	Handle(*Message) error
}

// MessageHandlerFunc is a convenience type to create MessageHandler implementations from a function
type MessageHandlerFunc func(*Message) error

// Handle is the MessageHandler interface implementation
func (m MessageHandlerFunc) Handle(msg *Message) error {
	return m(msg)
}
