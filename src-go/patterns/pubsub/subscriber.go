package pubsub

import (
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
)

// SubscribeFunc is called by a Subscriber when a message is received from a Publisher.
type SubscribeFunc func(*Message) error

// NewSubscriber returns a new Subscriber that will subscribe to topics published by an io.ReadWriter.
// The Connector may be a valid core.Dialer, or core.Listener.
func NewSubscriber(c core.Connector, topics ...string) *Subscriber {
	check.MustGreater(len(topics), 0, "number of topics")

	sub := &Subscriber{
		topics:    topics,
		connector: c,
		proto:     &subProtoV1{},
	}

	sub.connector.OnConnect(sub.onNewConn)
	sub.ch = make(chan Message, sub.connector.QueueSize())
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
	errFunc   core.ErrorFunc
	quit      chan bool
	finished  chan bool
}

// run is the engine of the Subscriber.
//
// TODO: handle connection errors by retrying.
func (sub *Subscriber) run() {
	defer func() {
		log.Println("subscriber finished.")
		sub.finished <- true
	}()

	check.Must(sub.connector.Open(sub.proto))

	for {
		select {
		case <-sub.quit:
			core.CloseQueues(sub.connector)
			return
		default:
			err := core.RecvQueues(sub.connector, sub.onRecv, sub.onRecvError)
			if err != nil {
				core.CloseQueues(sub.connector)
				return
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

// onRecv forwards a received message to the Subscribers SubscribeFunc if it's configured.
func (sub *Subscriber) onRecv(v interface{}) error {
	if sub.subFunc == nil {
		return nil
	}
	return sub.subFunc(v.(*Message))
}

// onRecvError forwards a Queue receive error to the Subscribers ErrorFunc if it's configured.
func (sub *Subscriber) onRecvError(err error) error {
	if sub.errFunc == nil {
		return nil
	}
	return sub.errFunc(err)
}

// Error is called to configure the ErrorFunc of a Subscriber,
// which is executed if a runtime error occurs while subscribing.
func (sub *Subscriber) Error(errFunc core.ErrorFunc) {
	sub.errFunc = errFunc
}

// Subscribe is called to configure the SubscribeFunc of a Subscriber,
// which is executed every time data is received from a Publisher.
//
// TODO: cater for multiple calls and from multiple goroutines.
// TODO: how to cater for Close() then Subscribe() ?
func (sub *Subscriber) Subscribe(subFunc SubscribeFunc) {
	check.MustNotNil(subFunc, "subscription function")

	sub.subFunc = subFunc
	go sub.run()
}

// Close is called to close any open connections and invalidate the Subscriber.
func (sub *Subscriber) Close() error {
	sub.quit <- true
	<-sub.finished
	return nil
}
