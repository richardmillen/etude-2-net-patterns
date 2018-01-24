package pubsub

import (
	"errors"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
)

// NewPublisher returns a new Publisher that will publish messages to Subscriber's.
// TODO: figure out good way to set queue size without cluttering the API.
func NewPublisher(c core.Connector) *Publisher {
	pub := &Publisher{connector: c}

	pub.ch = make(chan Message, 1)
	pub.quit = make(chan bool, 1)
	pub.finished = make(chan bool)

	go pub.run()
	return pub
}

// Publisher sends messages to zero or more Subscriber's.
type Publisher struct {
	connector core.Connector
	ch        chan Message
	quit      chan bool
	finished  chan bool
}

// run is the engine of the Publisher.
// TODO: should we report subscription queue errors to the consumer?
func (pub *Publisher) run() {
	defer func() {
		log.Println("publisher done.")
		pub.finished <- true
	}()

	check.Must(pub.connector.Open(&pubProtoV1{}))

	for {
		select {
		case <-pub.quit:
			return
		case m := <-pub.ch:
			queues := pub.connector.GetQueues()
			for _, q := range queues {
				err := q.Send(&m)
				check.Log(err)
			}
		}
	}
}

// Publish sends data to subscribers.
func (pub *Publisher) Publish(topic string, content []byte) error {
	select {
	case pub.ch <- Message{Topic: topic, Body: content}:
		return nil
	default:
		return errors.New("publisher queue full")
	}
}

// Close is called to stop and invalidate the Publisher.
func (pub *Publisher) Close() error {
	pub.quit <- true
	<-pub.finished
	return nil
}
