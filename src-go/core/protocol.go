package core

// QueueGreeter is the interface that wraps the protocol Greet method.
type QueueGreeter interface {
	Greet(q *Queue) error
}

// QueueSender is the interface that wraps the protocol Send method.
type QueueSender interface {
	Send(q *Queue, v interface{}) error
}

// QueueReceiver is the interface that wraps the protocol Recv method.
type QueueReceiver interface {
	Recv(q *Queue) (interface{}, error)
}

// SendReceiver is the interface that groups the Greet, Send and Recv methods.
type SendReceiver interface {
	QueueSender
	QueueReceiver
}

// GreetSendReceiver is the interface that groups the Greet, Send and Recv methods.
type GreetSendReceiver interface {
	QueueGreeter
	SendReceiver
}
