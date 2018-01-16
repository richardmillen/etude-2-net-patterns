package pubsub

import (
	"log"
	"net"
)

// NewPublisher returns a new Publisher that will publish messages to Subscriber's.
func NewPublisher(lnr net.Listener) *Publisher {
	pub := &Publisher{lnr: newListener(lnr)}

	pub.ch = make(chan data, 1)
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

func (pub *Publisher) run() {
	defer log.Println("publisher done.")

	for {
		select {
		case <-pub.quit:
			return
		case m := <-pub.ch:
			subs := pub.lnr.getSubscriptions()
			for _, s := range subs {
				s.receive(m)
			}
			break
		}
	}
}

// Publish sends data to subscribers.
func (pub *Publisher) Publish(topic string, body []byte) {
	pub.ch <- data{topic: topic, body: body}
}

// Close closes and invalidates the Publisher.
func (pub *Publisher) Close() error {
	close(pub.quit)
	return nil
}
