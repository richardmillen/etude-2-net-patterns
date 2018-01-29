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
	d := &Dialer{
		Network:    network,
		RemoteAddr: address,
		EP:         NewHostEndpoint(),
		queueSize:  DefQueueSize,
		queues:     make([]*Queue, 0, 1),
		Dialer: net.Dialer{
			Timeout: DefDialerTimeout,
		},
	}
	d.wgOpen.Add(1)
	return d
}

// A Dialer enables an endpoint to connect to another endpoint.
type Dialer struct {
	net.Dialer
	Network    string
	RemoteAddr string
	EP         *Endpoint
	queueSize  int
	queues     []*Queue
	conn       net.Conn
	connFunc   ConnectFunc
	gsr        GreetSendReceiver
	wgOpen     sync.WaitGroup
}

// QueueSize returns the size used when creating new connection Queues.
func (d *Dialer) QueueSize() int {
	return d.queueSize
}

// Open is called to initialise the Dialer with a protocol and connect.
//
// An id is often used so the remote Endpoint is able to uniquely identify this
// Connector (local Endpoint).
func (d *Dialer) Open(gsr GreetSendReceiver) (err error) {
	defer d.wgOpen.Done()

	d.gsr = gsr

	d.conn, err = d.Dial(d.Network, d.RemoteAddr)
	if err != nil {
		return
	}
	d.EP.Addr = GetEndpointAddress(d.conn.LocalAddr())

	d.queues = append(d.queues, NewQueue(d.conn, d.QueueSize()))
	d.queues[0].SetProp(PropAddressKey, d.EP.Addr)

	if d.connFunc != nil {
		err = d.connFunc(d.queues[0])
		if err != nil {
			return
		}
	}

	err = d.gsr.Greet(d.queues[0])
	if err != nil {
		return
	}

	err = check.NotNil(d.queues[0].sr, "queue send-receiver")
	if err != nil {
		return
	}

	return
}

// GetQueues returns a slice of active endpoint Queues.
//
// A Dialer currently only supports one queue
// (stored in a slice to avoid recreating slices unnecessarily),
// but it could be extended to support 'n' queues.
func (d *Dialer) GetQueues() []*Queue {
	d.wgOpen.Wait()
	return d.queues
}

// Close is called to close any open connections.
func (d *Dialer) Close() error {
	return nil
}

// OnConnect sets a ConnectFunc to be invoked whenever a new connection Queue is created.
func (d *Dialer) OnConnect(connFunc ConnectFunc) {
	d.connFunc = connFunc
}
