package pattern

import (
	"errors"
	"io"
)

// NewSubscriber returns a new Subscriber that will subscribe to an io.ReadWriter.
func NewSubscriber(rw io.ReadWriter) *Subscriber {
	return nil
}

// Publisher is used to publish a message stream to zero or more Subscriber's.
type Publisher struct {
}

// Subscriber is used to subscribe to a Publisher service.
type Subscriber struct {
}

// Subscribe is called to receive data from a publisher
func (s *Subscriber) Subscribe(fun func([]byte), filters ...string) error {
	return errors.New("not implemented")
}
