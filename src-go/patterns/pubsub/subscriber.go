package pubsub

import (
	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/core"
)

// SubscribeFunc is called by a Subscriber when a message is received from a Publisher.
type SubscribeFunc func(*Message) error

// NewSubscriber returns a new Subscriber that will subscribe to topics published by an io.ReadWriter.
// The Connector may be a valid core.Dialer, or core.Listener.
func NewSubscriber(c core.Connector, topics ...string) *Subscriber {
	check.MustGreater(len(topics), 0, "number of topics")

	sub := &Subscriber{
		topics:  topics,
		service: core.NewService(c, &subProtoV1{}),
	}
	sub.service.Connect(sub.onNewConn)
	sub.service.Error(sub.onError)
	sub.service.Recv(sub.onRecv)
	return sub
}

// Subscriber subscribes to topics published by a Publisher.
type Subscriber struct {
	topics  []string
	service *core.Service
	subFunc SubscribeFunc
	errFunc core.ErrorFunc
}

// onNewConn is invoked by the Subscribers Connector whenever a new connection Queue is created.
//
// TODO: cater for multiple topics.
func (sub *Subscriber) onNewConn(q *core.Queue) error {
	q.SetProp(propTopicKey, sub.topics[0])
	return nil
}

// onError forwards a Queue receive error to the Subscribers ErrorFunc if it's configured.
func (sub *Subscriber) onError(err error) error {
	if sub.errFunc == nil {
		return nil
	}
	return sub.errFunc(err)
}

// onRecv forwards a received message to the Subscribers SubscribeFunc if it's configured.
func (sub *Subscriber) onRecv(v interface{}) error {
	if sub.subFunc == nil {
		return nil
	}
	return sub.subFunc(v.(*Message))
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
	sub.subFunc = subFunc
}

// Start is called to begin the Subscriber.
func (sub *Subscriber) Start() {
	sub.service.Start()
}

// Close is called to close any open connections and invalidate the Subscriber.
func (sub *Subscriber) Close() error {
	return sub.service.Close()
}
