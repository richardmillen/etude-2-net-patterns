package pubsub

import (
	"errors"
	"io"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
)

// ConnectFunc is called by a Subscriber to connect to a Publisher.
type ConnectFunc func() (io.ReadWriteCloser, error)

// SubscribeFunc is called by a Subscriber when a message is received from a Publisher.
type SubscribeFunc func(*Message) error

// NewSubscriber returns a new Subscriber that will subscribe to topics published by an io.ReadWriter.
// The id is required so that the Publisher is able to uniquely identify the subscriber.
// The Connector may be a valid core.Dialer, or core.Listener.
func NewSubscriber(id string, c core.Connector) *Subscriber {
	sub := &Subscriber{id: id}
	sub.connector = c
	sub.proto = &subProtoV1{}
	sub.ch = make(chan Message, core.DefQueueSize)
	sub.quit = make(chan bool, 1)
	sub.finished = make(chan bool)
	return sub
}

// Subscriber subscribes to topics published by a Publisher.
type Subscriber struct {
	id        string
	ch        chan Message
	connector core.Connector
	proto     SubProtocol
	subFunc   SubscribeFunc
	topics    []string
	quit      chan bool
	finished  chan bool
}

// run is the engine of the Subscriber.
func (sub *Subscriber) run() {
	defer func() {
		log.Println("subscriber done.")
		sub.finished <- true
	}()

	check.Must(sub.connector.Open(&pubProtoV1{}))

	for {
		select {
		case <-sub.quit:
			return
		default:
			queues := sub.connector.GetQueues()
			for _, q := range queues {
				m, err := sub.proto.Recv(q)
				if !check.Log(err) {
					continue
				}

				check.Must(sub.subFunc(m.(*Message)))
			}
		}
	}
}

// Subscribe receives data from one or more publishers.
//
// TODO: cater for multiple calls and from multiple goroutines.
// TODO: how to cater for Close() then Subscribe() ?
func (sub *Subscriber) Subscribe(subFunc SubscribeFunc, topics ...string) (err error) {
	if sub.topics != nil {
		err = errors.New("already subscribing")
		return
	}

	sub.subFunc = subFunc
	sub.topics = make([]string, len(topics))
	copy(sub.topics, topics)

	go sub.run()

	return
}

// Close is called to close any open connections and invalidate the Subscriber.
func (sub *Subscriber) Close() error {
	defer func() {
		<-sub.finished
	}()

	sub.quit <- true
	return sub.connector.Close()
}
