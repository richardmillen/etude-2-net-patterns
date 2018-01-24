package pubsub

import (
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
)

// ConnectFunc is called by a Subscriber to connect to a Publisher.
//type ConnectFunc func() (io.ReadWriteCloser, error)

// SubscribeFunc is called by a Subscriber when a message is received from a Publisher.
type SubscribeFunc func(*Message) error

// ErrorFunc is called by a Subscriber when an error occurs during processing.
type ErrorFunc func(error)

// NewSubscriber returns a new Subscriber that will subscribe to topics published by an io.ReadWriter.
// The Connector may be a valid core.Dialer, or core.Listener.
func NewSubscriber(c core.Connector, topics ...string) *Subscriber {
	check.IsGreater(len(topics), 0, "number of topics")

	sub := &Subscriber{
		topics:    topics,
		connector: c,
		proto:     &subProtoV1{},
	}

	sub.connector.OnConnect(sub.onNewConn)
	sub.ch = make(chan Message, core.DefQueueSize)
	sub.quit = make(chan bool, 1)
	sub.finished = make(chan bool)
	return sub
}

// Subscriber subscribes to topics published by a Publisher.
type Subscriber struct {
	topics    []string
	ch        chan Message
	connector core.Connector
	proto     SubProtocol
	subFunc   SubscribeFunc
	errFunc   ErrorFunc
	quit      chan bool
	finished  chan bool
}

// run is the engine of the Subscriber.
//
// TODO: handle connection errors by retrying.
func (sub *Subscriber) run() {
	defer func() {
		log.Println("subscriber done.")
		sub.finished <- true
	}()

	check.Must(sub.connector.Open(sub.proto))

	for {
		select {
		case <-sub.quit:
			return
		default:
			queues := sub.connector.GetQueues()
			for _, q := range queues {
				m, err := sub.proto.Recv(q)
				err = patterns.Error(err)

				switch err.(type) {
				case patterns.ErrOffline:
					log.Printf("error: %s. aborting...\n", err)
					return
				case patterns.ErrConnLost:
					log.Printf("error: %s. aborting...", err)
					return
				case nil:
				default:
					check.Log(err)
					continue
				}

				if check.Log(sub.subFunc(m.(*Message))) {
					return
				}
			}
		}
	}
}

// onNewConn is invoked by the Subscribers Connector whenever a new connection Queue is created.
//
// TODO: cater for multiple topics.
func (sub *Subscriber) onNewConn(q *core.Queue) error {
	q.SetProp(propTopicKey, sub.topics[0])
	return nil
}

// Error receives a runtime error if one should occur while subscribing.
func (sub *Subscriber) Error(errFunc ErrorFunc) {
	sub.errFunc = errFunc
}

// Subscribe receives data from one or more publishers.
//
// TODO: cater for multiple calls and from multiple goroutines.
// TODO: how to cater for Close() then Subscribe() ?
func (sub *Subscriber) Subscribe(subFunc SubscribeFunc) {
	check.IsNotNil(subFunc, "subscription function")

	sub.subFunc = subFunc
	go sub.run()
}

// Close is called to close any open connections and invalidate the Subscriber.
func (sub *Subscriber) Close() error {
	sub.quit <- true
	<-sub.finished
	return nil
}
