package pubsub

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/pubsub/internal"
)

// newQueue constructs a new subscription queue (publisher-side connection).
func newQueue(conn net.Conn, queueSize uint, quit chan bool, wg *sync.WaitGroup) *queue {
	q := &queue{conn: conn, quit: quit, wg: wg}
	q.ch = make(chan internal.Msg, queueSize)
	go q.run()
	return q
}

// queue handles a subscriber connection on the Publisher.
type queue struct {
	conn  net.Conn
	proto PubProtocol
	id    string
	topic string
	ch    chan internal.Msg
	quit  chan bool
	wg    *sync.WaitGroup
}

func (q *queue) run() {
	defer q.wg.Done()

	for {
		select {
		case m := <-q.ch:
			if !q.subscribing(m.Topic) {
				break
			}

			err := q.proto.Send(q.conn, m.Topic, m.Body)
			if check.Log(err) {
				return
			}
		case <-q.quit:
			return
		}
	}
}

func (q *queue) subscribing(topic string) bool {
	return strings.HasPrefix(topic, q.topic)
}

func (q *queue) send(m internal.Msg) error {
	select {
	case q.ch <- m:
		return nil
	default:
		return fmt.Errorf("subscription queue '%s' full", q.id)
	}
}

// Conn returns a ReadWriteCloser that represents the connection to the Subscriber.
// Part of the Subscription interface.
func (q *queue) Conn() io.ReadWriteCloser {
	return q.conn
}

// SetProtocol is called to provide a queue with a version of the Pub protocol.
// Part of the Subscription interface.
func (q *queue) SetProtocol(p PubProtocol) {
	q.proto = p
}

// SetProps is called to set properties of the queue.
// Part of the Subscription interface.
//
// TODO: surely the queue shouldn't be reading properties that could become
// version-specific.
func (q *queue) SetProps(p map[string]string) {
	q.id = p[propIDKey]
	q.topic = p[propTopicKey]
}
