package pubsub

import (
	"github.com/richardmillen/etude-2-net-patterns/src-go/core"
)

// NewPublisher returns a new Publisher that will publish messages to Subscriber's.
func NewPublisher(c core.Connector) *Publisher {
	return &Publisher{service: core.NewService(c, &pubProtoV1{})}
}

// Publisher sends messages to zero or more Subscriber's.
type Publisher struct {
	service *core.Service
}

// Start is called to start the Publisher.
func (pub *Publisher) Start() {
	pub.service.Start()
}

// Publish sends data to subscribers.
func (pub *Publisher) Publish(topic string, content []byte) error {
	return pub.service.Send(&Message{Topic: topic, Body: content})
}

// Close is called to stop and invalidate the Publisher.
func (pub *Publisher) Close() error {
	return pub.service.Close()
}
