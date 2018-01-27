package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

// DefQueueSize is intended to be used as the default size for connection Queues.
const DefQueueSize = 100

const (
	// PropAddressKey is the key/name of the 'addr' Queue property.
	PropAddressKey = "addr"
)

// CloseConnectedQueues is called to close all connection Queues.
func CloseConnectedQueues(c Connector) {
	log.Println("core.CloseQueues: 1")

	queues := c.GetQueues()
	for _, q := range queues {
		q.Close()
	}

	log.Println("core.CloseQueues: 2")
}

// newQueue constructs a new connection queue.
func newQueue(conn net.Conn, queueSize uint) *Queue {
	q := &Queue{
		id:   uuid.New(),
		conn: conn,
	}
	q.props = make(map[string]interface{})
	q.ch = make(chan interface{}, queueSize)
	q.err = make(chan error, queueSize)
	q.quit = make(chan bool, 1)
	go q.run()
	return q
}

// Queue handles a subscriber connection on the Publisher.
type Queue struct {
	id         uuid.Bytes
	conn       net.Conn
	proto      StreamProtocol
	props      map[string]interface{}
	ch         chan interface{}
	err        chan error
	quit       chan bool
	wgSend     sync.WaitGroup
	wgFinished sync.WaitGroup
	//busy       bool
	//busyMutex  sync.Mutex
}

// ID is the unique identifier of the queue
//
// TODO: is this even needed?
func (q *Queue) ID() uuid.Bytes {
	return q.id
}

// Conn returns a ReadWriteCloser that represents the connection to the Subscriber.
func (q *Queue) Conn() io.ReadWriteCloser {
	return q.conn
}

// SetProtocol is called to provide a queue with a version of the Pub protocol.
func (q *Queue) SetProtocol(proto StreamProtocol) {
	q.proto = proto
}

// Prop returns the value of a property on the queue.
func (q *Queue) Prop(key string) interface{} {
	return q.props[key]
}

// SetProp is called to set a property on the queue.
func (q *Queue) SetProp(key string, value interface{}) {
	q.props[key] = value
}

// Send is called to pass data to a connection Queue.
func (q *Queue) Send(v interface{}) error {
	log.Println("Queue.Send: enter.")
	q.wgSend.Add(1)

	select {
	case q.ch <- v:
		log.Println("Queue.Send: message queued.")
		return nil
	default:
		q.wgSend.Done()
		return fmt.Errorf("connection queue '%s' is full", q.id)
	}
}

// Close waits for the Queue to finish processing messages then kills it.
func (q *Queue) Close() error {
	q.wgSend.Wait()
	q.quit <- true
	return nil
}

// run is the engine of a connection Queue.
func (q *Queue) run() {
	defer q.wgFinished.Done()

	for {
		select {
		case v := <-q.ch:
			log.Println("queue, message received:", v)
			err := q.proto.Send(q, v)
			if check.Log(err) {
				q.err <- err
			}
			q.wgSend.Done()
		case <-q.quit:
			log.Println("Queue.run: quit received.")
			return
		}
	}
}
