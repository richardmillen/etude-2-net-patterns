package core

import (
	"net"
	"sync"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// DefDialerTimeout is the default timeout value used by a Dialer.
const DefDialerTimeout = time.Second

// NewDialer constructs a new connection Dialer.
func NewDialer(network, address string) *Dialer {
	return &Dialer{
		Network: network,
		Address: address,
		Dialer: net.Dialer{
			Timeout: DefDialerTimeout,
		},
		q: make([]*Queue, 1),
	}
}

// A Dialer enables an endpoint to connect to another endpoint.
type Dialer struct {
	net.Dialer
	Network   string
	Address   string
	QueueSize uint
	proto     StreamProtocol
	q         []*Queue
	quit      chan bool
	wg        sync.WaitGroup
}

// Open is called to initialise the Dialer with a protocol.
func (d *Dialer) Open(proto StreamProtocol) error {
	d.proto = proto

	conn, err := d.Dial(d.Network, d.Address)
	if err != nil {
		return err
	}

	d.wg.Add(1)
	d.q[0] = newQueue(conn, d.QueueSize, d.quit, &(d.wg))

	go func() {
		err = d.proto.Greet(d.q[0])
		if check.Log(err) {
			return
		}
	}()

	return nil
}

// GetQueues returns a slice of active endpoint Queues.
//
// A Dialer currently only supports one queue
// (stored in a slice to avoid recreating slices unnecessarily),
// but it could be extended to support 'n' queues.
func (d *Dialer) GetQueues() []*Queue {
	return d.q
}

// Close is called to close any open connections.
func (d *Dialer) Close() error {
	close(d.quit)
	d.wg.Wait()
	return nil
}
