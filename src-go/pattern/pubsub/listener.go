package pubsub

import (
	"log"
	"net"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
)

func newListener(l net.Listener) listener {
	listener := &listener{inner: l}
	listener.receivers = make(map[string]*receiver)

	go listener.listen()

	return listener
}

// listener manages subscriber connections (receivers) for a Publisher.
type listener struct {
	inner     net.Listener
	receivers map[string]*receiver
	m         sync.Mutex
	wg        sync.WaitGroup
}

func (l *listener) listen() {
	defer log.Println("listener done.")

	for {
		conn, err := l.inner.Accept()
		// TODO: check the error to see if we've been closed!
		utils.CheckError(err)

		l.wg.Add(1)
		r := newReceiver(conn, l.wg)

		l.m.Lock()
		l.receivers[r.id] = r
		l.m.Unlock()
	}
}

func (l *listener) getReceivers() []*receiver {
	l.m.Lock()

	receivers := make([]*receiver, len(l.receivers))
	n := 0
	for _, r := range l.receivers {
		receivers[n] = r
		n++
	}

	l.m.Unlock()

	return receivers
}

func (l *listener) Close() error {
	log.Println("listener closing...")
	defer func() {
		l.wg.Wait()
		log.Println("listener closed.")
	}()

	return l.inner.Close()
}
