package core

// Connector represents an endpoint connection mechanism.
type Connector interface {
	GetQueues() []*Queue
	Open(proto StreamProtocol) error
	Close() error
}
