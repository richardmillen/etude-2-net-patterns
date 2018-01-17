package pubsub

import "io"

// ProtocolSignature is used to identify some messages belonging to the Pub-Sub protocol.
var ProtocolSignature = [...]byte{0x01, 0x01}

const (
	// PropIDKey is the identifier of the 'id' message property.
	propIDKey = "id"
	// PropTopicKey is the identifier of the 'topic' message property.
	propTopicKey = "topic"
)

// SubscriptionQueue is the interface of a subscription queue.
type SubscriptionQueue interface {
	Conn() io.ReadWriteCloser
	SetProtocol(p PubProtocol)
	SetProps(p map[string]string)
}

// PubProtocol is the Publisher-side of the pub-sub wire protocol.
type PubProtocol interface {
	Greet(q SubscriptionQueue) error
	Send(conn io.ReadWriter, topic string, body []byte) error
}

// SubProtocol is the Subscriber-side of the pub-sub wire protocol.
type SubProtocol interface {
	Ready(sub *Subscriber) error
	Recv(conn io.ReadWriter) (topic string, body []byte, err error)
}

// Greeting is a message sent by the publisher immediately after a subscriber connects.
type Greeting struct {
	Signature [2]byte
	Major     uint8
	Minor     uint8
}

// Ready is a message sent by the subscriber in response to a greeting.
type Ready struct {
	Major uint8
	Minor uint8
}
