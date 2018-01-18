package pubsub

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// newQueue constructs a new subscription queue (publisher-side connection).
func newQueue(conn net.Conn, queueSize uint, quit chan bool, wg *sync.WaitGroup) *Queue {
	q := &Queue{conn: conn, quit: quit, wg: wg}
	q.ch = make(chan Message, queueSize)
	go q.run()
	return q
}

// Queue handles a subscriber connection on the Publisher.
type Queue struct {
	conn  net.Conn
	proto PubProtocol
	id    string
	topic string
	ch    chan Message
	quit  chan bool
	wg    *sync.WaitGroup
}

// Conn returns a ReadWriteCloser that represents the connection to the Subscriber.
func (q *Queue) Conn() io.ReadWriteCloser {
	return q.conn
}

// SetProtocol is called to provide a queue with a version of the Pub protocol.
func (q *Queue) SetProtocol(p PubProtocol) {
	q.proto = p
}

// SetProps is called to set properties of the queue.
//
// TODO: surely the queue shouldn't be reading properties that could become
// version-specific.
func (q *Queue) SetProps(p map[string]string) {
	q.id = p[propIDKey]
	q.topic = p[propTopicKey]
}

func (q *Queue) run() {
	defer q.wg.Done()

	for {
		select {
		case <-q.quit:
			return
		case m := <-q.ch:
			if !q.subscribing(m.Topic) {
				break
			}

			err := q.proto.Send(q.conn, &m)
			if check.Log(err) {
				return
			}
		}
	}
}

func (q *Queue) subscribing(topic string) bool {
	return strings.HasPrefix(topic, q.topic)
}

func (q *Queue) send(m *Message) error {
	select {
	case q.ch <- *m:
		return nil
	default:
		return fmt.Errorf("subscription queue '%s' is full", q.id)
	}
}
