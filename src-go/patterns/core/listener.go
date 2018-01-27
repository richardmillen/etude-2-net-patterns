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
	proto     StreamProtocol
	queueSize uint
	queues    []*Queue
	connFunc  ConnectFunc
	m         sync.Mutex
	finished  chan bool
}

// Open is called to start the listener.
func (l *Listener) Open(proto StreamProtocol) error {
	l.proto = proto

	go func() {
		defer func() {
			log.Println("listener done.")
			l.finished <- true
		}()

		for {
			conn, err := l.Listener.Accept()
			if check.Log(err) {
				break
			}

			go func() {
				q := newQueue(conn, l.queueSize)

				if l.connFunc != nil {
					err = l.connFunc(q)
					if check.Log(err) {
						return
					}
				}

				err = l.proto.Greet(q)
				if check.Log(err) {
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
		select {
		case e := <-q.err:
			log.Printf("error reported by queue '%s': %s\n", q.ID(), e)
		default:
			qs = append(qs, q)
		}
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
