package core

import (
	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// ErrorFunc is callback used to report errors to a consumer.
// The return error is intended to be used to abort processing.
type ErrorFunc func(error) error

// RecvFunc is called by RecvAllQueues when data is received on a Queue.
type RecvFunc func(v interface{}) error

// ConnectFunc should be called whenever a new connection Queue is created.
type ConnectFunc func(q *Queue) error

// Connector represents an endpoint connection mechanism.
//
// TODO: create new Queues type and turn ...Queues functions into methods on type(?)
type Connector interface {
	QueueSize() int
	GetQueues() []*Queue
	Open(gsr GreetSendReceiver) error
	Close() error
	OnConnect(connFunc ConnectFunc)
}

// SendToQueues is called to send a message to all active Queues.
// n.b. Any send errors are held against each Queue.
func SendToQueues(c Connector, v interface{}) {
	queues := c.GetQueues()
	for _, q := range queues {
		err := q.Send(v)
		check.Log(err)
	}
}

// RecvQueues is called to receive on all Queues.
func RecvQueues(c Connector, recvFunc RecvFunc, errFunc ErrorFunc) (err error) {
	queues := c.GetQueues()
	for _, q := range queues {
		var m interface{}

		m, err = q.Recv()
		if err != nil {
			err = errFunc(err)
			if err != nil {
				return
			}
			continue
		}

		err = recvFunc(m)
		if err != nil {
			return
		}
	}
	return
}

// CloseQueues is called to close all connection Queues.
func CloseQueues(c Connector) {
	queues := c.GetQueues()
	for _, q := range queues {
		q.Close()
	}
}
