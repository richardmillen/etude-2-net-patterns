package pubsub

import (
	"log"
	"net"
)

// NewPublisher returns a new Publisher that will publish messages to Subscriber's.
func NewPublisher(l net.Listener) *Publisher {
	p := &Publisher{l: newListener(l)}

	p.ch = make(chan message, 1)
	p.quit = make(chan bool)

	go p.run()
	return p
}

// Publisher is used to publish a message stream to zero or more Subscriber's.
type Publisher struct {
	l    listener
	ch   chan message
	quit chan bool
}

func (p *Publisher) run() {
	defer log.Println("publisher done.")

	for {
		select {
		case m := <-p.ch:
			receivers := p.l.getReceivers()
			for _, r := range receivers {
				r.receive(m)
			}
			break
		case <-p.quit:
			return
		}
	}
}

// Publish sends data to subscribers.
func (p *Publisher) Publish(topic string, body []byte) {
	p.ch <- message{topic: topic, body: body}
}

// Close closes and invalidates the Publisher.
func (p *Publisher) Close() error {
	close(p.quit)
	return nil
}
