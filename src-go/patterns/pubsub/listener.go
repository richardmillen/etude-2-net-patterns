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
	lnr.subs = make(map[string]*subscription)
	lnr.quit = make(chan bool)
	lnr.proto = &protoV1{}

	go lnr.listen()
	return lnr
}

// listener manages publisher-side connections.
type listener struct {
	inner     net.Listener
	proto     protocol
	queueSize uint
	subs      map[string]*subscription
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
			sub := newSubscription(conn, lnr.queueSize, lnr.quit, &(lnr.wg))

			err = lnr.proto.greet(sub)
			if check.Log(err) {
				return
			}

			lnr.m.Lock()
			lnr.subs[sub.id] = sub
			lnr.m.Unlock()
		}()
	}
}

func (lnr *listener) getSubscriptions() []*subscription {
	lnr.m.Lock()

	subs := make([]*subscription, len(lnr.subs))
	n := 0
	for _, r := range lnr.subs {
		subs[n] = r
		n++
	}

	lnr.m.Unlock()

	return subs
}

func (lnr *listener) Close() error {
	log.Println("listener closing...")
	defer func() {
		lnr.wg.Wait()
		log.Println("listener closed.")
	}()

	return lnr.inner.Close()
}
