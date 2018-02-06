package core

import (
	"log"
	"net"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// ListenTCP constructs a new Listener for a network endpoint.
func ListenTCP(network string, laddr *net.TCPAddr) (l *Listener, err error) {
	inner, err := net.ListenTCP(network, laddr)
	if err != nil {
		return
	}

	return &Listener{
		Listener:  inner,
		EP:        NewEndpoint(laddr.String()),
		queueSize: DefQueueSize,
		queues:    []*Queue{},
		finished:  make(chan bool),
	}, nil
}

// Listener manages publisher-side connections.
type Listener struct {
	net.Listener
	EP        *Endpoint
	gsr       GreetSendReceiver
	queueSize int
	queues    []*Queue
	connFunc  ConnectFunc
	m         sync.Mutex
	finished  chan bool
}

// QueueSize returns the size used when creating new connection Queues.
func (l *Listener) QueueSize() int {
	return l.queueSize
}

// Open is called to start the listener.
func (l *Listener) Open(gsr GreetSendReceiver) error {
	l.gsr = gsr

	go func() {
		defer func() {
			l.finished <- true
		}()

		for {
			conn, err := l.Listener.Accept()
			if check.Log(err) {
				break
			}

			go func() {
				q := NewQueue(conn, l.queueSize)

				if l.connFunc != nil {
					err = l.connFunc(q)
					if check.Log(err) {
						return
					}
				}

				if check.Log(l.gsr.Greet(q)) {
					return
				}

				if check.Log(check.NotNil(q.sr, "queue send-receiver")) {
					return
				}

				l.m.Lock()
				l.queues = append(l.queues, q)
				l.m.Unlock()
			}()
		}
	}()

	return nil
}

// GetQueues returns a slice of currently active queues.
//
// TODO: write an in-place deletion.
func (l *Listener) GetQueues() []*Queue {
	l.m.Lock()

	qs := make([]*Queue, 0, len(l.queues))
	for _, q := range l.queues {
		err := q.Err()
		if err != nil {
			log.Printf("removing queue '%s': %s\n", q.ID(), err)
			continue
		}
		qs = append(qs, q)
	}
	l.queues = qs

	l.m.Unlock()

	return l.queues
}

// Close is called to close all open connections and stop listening.
func (l *Listener) Close() (err error) {
	err = l.Listener.Close()
	<-l.finished
	return
}

// OnConnect sets a ConnectFunc to be invoked whenever a new connection Queue is created.
func (l *Listener) OnConnect(connFunc ConnectFunc) {
	l.connFunc = connFunc
}
