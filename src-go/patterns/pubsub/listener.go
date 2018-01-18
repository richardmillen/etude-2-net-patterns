package pubsub

import (
	"log"
	"net"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// newListener constructs a new listener for a Publisher.
func newListener(l net.Listener, queueSize uint) *listener {
	lnr := &listener{inner: l}
	lnr.queueSize = queueSize
	lnr.qs = make(map[string]*Queue)
	lnr.quit = make(chan bool)
	lnr.proto = &pubProtoV1{}

	go lnr.listen()
	return lnr
}

// listener manages publisher-side connections.
type listener struct {
	inner     net.Listener
	proto     PubProtocol
	queueSize uint
	qs        map[string]*Queue
	m         sync.Mutex
	quit      chan bool
	wg        sync.WaitGroup
}

func (lnr *listener) listen() {
	defer log.Println("listener done.")

	for {
		conn, err := lnr.inner.Accept()
		// TODO: check the error to see if we've been closed!
		check.Error(err)

		go func() {
			lnr.wg.Add(1)
			q := newQueue(conn, lnr.queueSize, lnr.quit, &(lnr.wg))

			err = lnr.proto.Greet(q)
			if check.Log(err) {
				return
			}

			lnr.m.Lock()
			lnr.qs[q.id] = q
			lnr.m.Unlock()
		}()
	}
}

func (lnr *listener) getQueues() []*Queue {
	lnr.m.Lock()

	qs := make([]*Queue, len(lnr.qs))
	n := 0
	for _, r := range lnr.qs {
		qs[n] = r
		n++
	}

	lnr.m.Unlock()

	return qs
}

func (lnr *listener) Close() error {
	log.Println("listener closing...")
	defer func() {
		lnr.wg.Wait()
		log.Println("listener closed.")
	}()

	close(lnr.quit)
	return lnr.inner.Close()
}
