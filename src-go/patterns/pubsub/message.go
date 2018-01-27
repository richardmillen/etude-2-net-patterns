package pubsub

import "fmt"

// Message is a message from a Publisher.
type Message struct {
	Topic string
	Body  []byte
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %s", m.Topic, string(m.Body))
}
