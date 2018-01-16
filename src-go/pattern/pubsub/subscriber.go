package pubsub

import (
	"errors"
	"io"
)

// NewSubscriber returns a new Subscriber that will subscribe to an io.ReadWriter.
func NewSubscriber(rw io.ReadWriter) *Subscriber {
	return nil
}

// Subscriber is used to subscribe to a Publisher service.
type Subscriber struct {
}

// Subscribe receives data from one or more publishers.
func (s *Subscriber) Subscribe(fun func([]byte), filters ...string) error {
	return errors.New("not implemented")
}
