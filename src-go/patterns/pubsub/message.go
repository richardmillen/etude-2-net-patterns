package pubsub

// Message is a message from a Publisher.
type Message struct {
	Topic string
	Body  []byte
}
