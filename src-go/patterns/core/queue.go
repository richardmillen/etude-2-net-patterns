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

// newQueue constructs a new subscription queue (publisher-side connection).
func newQueue(conn net.Conn, queueSize uint, quit chan bool, wg *sync.WaitGroup) *Queue {
	q := &Queue{
		id:   uuid.New(),
		conn: conn,
		quit: quit,
		wg:   wg,
	}
	q.props = make(map[string]interface{})
	q.ch = make(chan interface{}, queueSize)
	q.err = make(chan error)
	go q.run()
	return q
}

// Queue handles a subscriber connection on the Publisher.
type Queue struct {
	id    uuid.Bytes
	conn  net.Conn
	proto StreamProtocol
	props map[string]interface{}
	ch    chan interface{}
	err   chan error
	quit  chan bool
	wg    *sync.WaitGroup
	m     sync.Mutex
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

func (q *Queue) run() {
	defer q.wg.Done()

	for {
		select {
		case v := <-q.ch:
			log.Println("queue, message received:", v)
			err := q.proto.Send(q, v)
			if check.Log(err) {
				q.err <- err
				return
			}
		default:
			select {
			case <-q.quit:
				return
			}
		}
	}
}

// Send is called to pass data to a connection Queue.
func (q *Queue) Send(v interface{}) error {
	select {
	case q.ch <- v:
		return nil
	default:
		return fmt.Errorf("connection queue '%s' is full", q.id)
	}
}
