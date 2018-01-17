package pubsub

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// newSubscription constructs a new subscription (publisher-side connection).
func newSubscription(conn net.Conn, queueSize uint, quit chan bool, wg *sync.WaitGroup) *subscription {
	s := &subscription{conn: conn, quit: quit, wg: wg}
	s.ch = make(chan data, queueSize)
	go s.enable()
	return s
}

// subscription handles a subscriber connection on the Publisher.
type subscription struct {
	conn  net.Conn
	proto protocol
	id    string
	topic string
	ch    chan data
	quit  chan bool
	wg    *sync.WaitGroup
}

func (sub *subscription) enable() {
	defer sub.wg.Done()

	for {
		select {
		case d := <-sub.ch:
			if !strings.HasPrefix(d.topic, sub.topic) {
				break
			}

			err := sub.proto.sendTo(sub, d.topic, d.content)
			if check.Log(err) {
				return
			}
		case <-sub.quit:
			return
		}
	}
}

func (sub *subscription) send(d data) error {
	select {
	case sub.ch <- d:
		return nil
	default:
		return fmt.Errorf("subscription '%s' queue full", sub.id)
	}
}
