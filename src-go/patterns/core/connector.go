package core

// ConnectFunc should be called whenever a new connection Queue is created.
type ConnectFunc func(q *Queue) error

// Connector represents an endpoint connection mechanism.
type Connector interface {
	GetQueues() []*Queue
	Open(proto StreamProtocol) error
	Close() error
	OnConnect(connFunc ConnectFunc)
}
