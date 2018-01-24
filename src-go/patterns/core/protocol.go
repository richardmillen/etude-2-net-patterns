package core

// StreamProtocol is intended to be used as a bare bones
// interface of a stream-oriented network wire protocol.
type StreamProtocol interface {
	Greet(q *Queue) error
	Send(q *Queue, v interface{}) error
	Recv(q *Queue) (interface{}, error)
}
