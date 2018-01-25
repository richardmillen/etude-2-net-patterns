package core

import (
	"log"
	"net"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
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
		queues:    make(map[string]*Queue),
		quit:      make(chan bool),
	}, nil
}

// Listener manages publisher-side connections.
type Listener struct {
	net.Listener
	EP        *Endpoint
	proto     StreamProtocol
	queueSize uint
	queues    map[string]*Queue
	connFunc  ConnectFunc
	m         sync.Mutex
	quit      chan bool
	wg        sync.WaitGroup
}

// Open is called to start the listener.
func (l *Listener) Open(proto StreamProtocol) error {
	l.proto = proto

	go func() {
		defer log.Println("listener done.")

		for {
			log.Println("waiting for subscriber connections...")

			conn, err := l.Accept()
			// TODO: check the error to see if we've been closed!
			check.Error(err)

			go func() {
				l.wg.Add(1)
				q := newQueue(conn, l.queueSize, l.quit, &(l.wg))

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
				l.queues[q.Prop(PropUUIDKey).(uuid.Bytes).String()] = q
				l.m.Unlock()
			}()
		}
	}()

	return nil
}

// GetQueues returns a slice of currently active queues.
func (l *Listener) GetQueues() []*Queue {
	l.m.Lock()

	qs := make([]*Queue, 0, len(l.queues))
	for _, q := range l.queues {
		select {
		case e := <-q.err:
			log.Printf("error reported by queue '%s': %s\n", q.Prop(PropUUIDKey).(uuid.Bytes), e)
			delete(l.queues, q.Prop(PropUUIDKey).(uuid.Bytes).String())
		default:
			qs = append(qs, q)
		}
	}

	l.m.Unlock()

	return qs
}

// Close is called to close all open connections and stop listening.
func (l *Listener) Close() error {
	defer func() {
		l.wg.Wait()
	}()

	close(l.quit)
	return l.Listener.Close()
}

// OnConnect sets a ConnectFunc to be invoked whenever a new connection Queue is created.
func (l *Listener) OnConnect(connFunc ConnectFunc) {

}