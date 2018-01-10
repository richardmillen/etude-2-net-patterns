package pattern

import (
	"errors"
	"io"
	"net"
)

// NewPublisher returns a new Publisher that will publish messages to Subscriber's.
func NewPublisher(listener net.Listener) *Publisher {
	return nil
}

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

// Publish sends data to subscribers.
func (p *Publisher) Publish(filter string, b []byte) {

}

// Close closes and invalidates the Publisher.
func (p *Publisher) Close() {

}

// Subscribe receives data from one or more publishers.
func (s *Subscriber) Subscribe(fun func([]byte), filters ...string) error {
	return errors.New("not implemented")
}
