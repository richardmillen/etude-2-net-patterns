package pubsub

import (
	"errors"
	"log"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

const outQueueSize = 10

// NewPublisher returns a new Publisher that will publish messages to Subscriber's.
// TODO: figure out good way to set queue size without cluttering the API.
func NewPublisher(lnr net.Listener) *Publisher {
	pub := &Publisher{lnr: newListener(lnr, outQueueSize)}

	pub.ch = make(chan Message, 1)
	pub.quit = make(chan bool)

	go pub.run()
	return pub
}

// Publisher sends messages to zero or more Subscriber's.
type Publisher struct {
	lnr  *listener
	ch   chan Message
	quit chan bool
}

// run is the engine of the Publisher.
// TODO: should we report subscription queue errors to the consumer?
func (pub *Publisher) run() {
	defer log.Println("publisher done.")

	for {
		select {
		case <-pub.quit:
			return
		case d := <-pub.ch:
			qs := pub.lnr.getQueues()
			for _, q := range qs {
				err := q.send(d)
				check.Log(err)
			}
		}
	}
}

// QueueSize returns the number of items that can be stored
// in each subscription queue.
func (pub *Publisher) QueueSize() uint {
	return pub.lnr.queueSize
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

// Close closes and invalidates the Publisher.
func (pub *Publisher) Close() error {
	pub.quit <- true
	return nil
}
