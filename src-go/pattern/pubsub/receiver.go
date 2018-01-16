package pubsub

import (
	"net"
	"sync"
)

func newReceiver(conn net.Conn, quit chan bool, wg *sync.WaitGroup) *receiver {
	r := &receiver{conn: conn, quit: quit, wg: wg}
	r.ch = make(chan message, 1)
	go r.run()
	return r
}

// receiver handles a subscriber connection.
type receiver struct {
	conn net.Conn
	id   string
	addr string
	ch   chan message
	quit chan bool
	wg   *sync.WaitGroup
}

func (r *receiver) run() {
	defer r.wg.Done()

	for {
		select {
		case m := <-r.ch:
			break
		case <-r.quit:
			return
		}
	}
}

func (r *receiver) receive(m message) {
	r.ch <- m
}
