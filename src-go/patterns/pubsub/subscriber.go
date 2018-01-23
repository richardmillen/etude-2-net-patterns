package pubsub

import (
	"errors"
	"io"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

const inQueueSize = 10

// ConnectFunc is called by a Subscriber to connect to a Publisher.
type ConnectFunc func() (io.ReadWriteCloser, error)

// SubscribeFunc is called by a Subscriber when a message is received from a Publisher.
type SubscribeFunc func(*Message) error

// NewSubscriber returns a new Subscriber that will subscribe to topics
// published by an io.ReadWriter.
func NewSubscriber(id string) *Subscriber {
	sub := &Subscriber{id: id}
	sub.ch = make(chan Message, inQueueSize)
	sub.quit = make(chan bool)
	sub.proto = &subProtoV1{}
	return sub
}

// Subscriber subscribes to topics published by a Publisher.
type Subscriber struct {
	id       string
	ch       chan Message
	proto    SubProtocol
	conn     io.ReadWriteCloser
	connFunc ConnectFunc
	subFunc  SubscribeFunc
	quit     chan bool
	topics   []string
}

// run is the engine of the Subscriber.
func (sub *Subscriber) run() {
	defer log.Println("subscriber done.")

	check.Must(sub.proto.Ready(sub))

	for {
		select {
		case <-sub.quit:
			return
		default:
			m, err := sub.proto.Recv(sub.conn)
			check.Error(err)

			check.Must(sub.subFunc(m))
		}
	}
}

// Subscribe receives data from one or more publishers.
// TODO: cater for multiple calls and from multiple goroutines.
// TODO: how to cater for Close() then Subscribe() ?
func (sub *Subscriber) Subscribe(connFunc ConnectFunc, subFunc SubscribeFunc, topics ...string) (err error) {
	if sub.topics != nil {
		err = errors.New("already subscribing")
		return
	}

	sub.topics = make([]string, len(topics))
	copy(sub.topics, topics)

	sub.connFunc = connFunc
	sub.subFunc = subFunc

	sub.conn, err = sub.connFunc()
	if err != nil {
		return
	}

	go sub.run()

	return
}

// Close is called to close any open connections to a Publisher.
func (sub *Subscriber) Close() (err error) {
	if sub.conn != nil {
		select {
		case sub.quit <- true:
		default:
		}

		err = sub.conn.Close()
		sub.conn = nil
	}
	return
}
