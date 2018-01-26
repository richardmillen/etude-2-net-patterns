package core

import (
	"net"
	"sync"
	"time"
)

// DefDialerTimeout is the default timeout value used by a Dialer.
const DefDialerTimeout = time.Second

// NewDialer constructs a new connection Dialer.
func NewDialer(network, address string) *Dialer {
	return &Dialer{
		Network:    network,
		RemoteAddr: address,
		EP:         NewHostEndpoint(),
		q:          make([]*Queue, 1),
		quit:       make(chan bool),
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
	QueueSize  uint
	q          []*Queue
	connFunc   ConnectFunc
	proto      StreamProtocol
	quit       chan bool
	wg         sync.WaitGroup
}

// Open is called to initialise the Dialer with a protocol and connect.
//
// An id is often used so the remote Endpoint is able to uniquely identify this
// Connector (local Endpoint).
func (d *Dialer) Open(proto StreamProtocol) error {
	d.proto = proto

	conn, err := d.Dial(d.Network, d.RemoteAddr)
	if err != nil {
		return err
	}
	d.EP.Addr = GetEndpointAddress(conn.LocalAddr())

	d.wg.Add(1)
	d.q[0] = newQueue(conn, d.QueueSize, d.quit, &(d.wg))
	d.q[0].SetProp(PropAddressKey, d.EP.Addr)

	if d.connFunc != nil {
		err = d.connFunc(d.q[0])
		if err != nil {
			return err
		}
	}

	err = d.proto.Greet(d.q[0])
	return err
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

// OnConnect sets a ConnectFunc to be invoked whenever a new connection Queue is created.
func (d *Dialer) OnConnect(connFunc ConnectFunc) {
	d.connFunc = connFunc
}
