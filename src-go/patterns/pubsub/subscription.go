package pubsub

import (
	"net"
	"strings"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
)

// newSubscription constructs a new subscription (publisher-side connection).
func newSubscription(conn net.Conn, quit chan bool, wg *sync.WaitGroup) *subscription {
	s := &subscription{conn: conn, quit: quit, wg: wg}
	s.ch = make(chan data, 1)
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

			err := sub.proto.send(d.topic, d.body)
			if utils.LogError(err) {
				return
			}
			break
		case <-sub.quit:
			return
		}
	}
}

func (sub *subscription) receive(d data) {
	sub.ch <- d
}
