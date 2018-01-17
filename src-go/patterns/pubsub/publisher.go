package pubsub

import (
	"errors"
	"log"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

const (
	pubQueueSize     = 2
	defaultQueueSize = 10
)

// NewPublisher returns a new Publisher that will publish messages to Subscriber's.
// TODO: figure out good way to set queue size without cluttering the API.
func NewPublisher(lnr net.Listener) *Publisher {
	pub := &Publisher{lnr: newListener(lnr, defaultQueueSize)}

	pub.ch = make(chan data, pubQueueSize)
	pub.quit = make(chan bool)

	go pub.run()
	return pub
}

// Publisher sends messages to zero or more Subscriber's.
type Publisher struct {
	lnr  *listener
	ch   chan data
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
			subs := pub.lnr.getSubscriptions()
			for _, sub := range subs {
				err := sub.send(d)
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
	case pub.ch <- data{topic: topic, content: content}:
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
