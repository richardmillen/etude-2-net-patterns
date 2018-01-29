package core

import (
	"log"
	"net"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// DefDialerTimeout is the default timeout value used by a Dialer.
const DefDialerTimeout = time.Second

// NewDialer constructs a new connection Dialer.
func NewDialer(network, address string) *Dialer {
	return &Dialer{
		Network:    network,
		RemoteAddr: address,
		EP:         NewHostEndpoint(),
		QueueSize:  DefQueueSize,
		q:          make([]*Queue, 1),
		Dialer: net.Dialer{
			Timeout: DefDialerTimeout,
		},
	}
}

// A Dialer enables an endpoint to connect to another endpoint.
type Dialer struct {
	net.Dialer
	Network    string
	RemoteAddr string
	EP         *Endpoint
	QueueSize  int
	q          []*Queue
	conn       net.Conn
	connFunc   ConnectFunc
	gsr        GreetSendReceiver
}

// Open is called to initialise the Dialer with a protocol and connect.
//
// An id is often used so the remote Endpoint is able to uniquely identify this
// Connector (local Endpoint).
func (d *Dialer) Open(gsr GreetSendReceiver) (err error) {
	d.gsr = gsr

	d.conn, err = d.Dial(d.Network, d.RemoteAddr)
	if err != nil {
		return
	}
	d.EP.Addr = GetEndpointAddress(d.conn.LocalAddr())

	log.Println("Dialer.Open: queue size:", d.QueueSize)
	d.q[0] = NewQueue(d.conn, d.QueueSize)
	d.q[0].SetProp(PropAddressKey, d.EP.Addr)

	if d.connFunc != nil {
		err = d.connFunc(d.q[0])
		if err != nil {
			return
		}
	}

	err = d.gsr.Greet(d.q[0])
	if err != nil {
		return
	}

	err = check.NotNil(d.q[0].sr, "queue send-receiver")
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
	return d.q
}

// Close is called to close any open connections.
func (d *Dialer) Close() error {
	return nil
}

// OnConnect sets a ConnectFunc to be invoked whenever a new connection Queue is created.
func (d *Dialer) OnConnect(connFunc ConnectFunc) {
	d.connFunc = connFunc
}
