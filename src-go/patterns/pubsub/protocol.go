package pubsub

import (
	"errors"
	"io"
)

// ProtocolSignature is used to identify some messages belonging to the Pub-Sub protocol.
var ProtocolSignature = [...]byte{0x01, 0x01}

const (
	// PropIDKey is the identifier of the 'id' message property.
	propIDKey = "id"
	// PropTopicKey is the identifier of the 'topic' message property.
	propTopicKey = "topic"
)

// PubProtocol is the Publisher-side of the pub-sub wire protocol.
type PubProtocol interface {
	Greet(q *Queue) error
	Send(w io.Writer, m *Message) error
}

// SubProtocol is the Subscriber-side of the pub-sub wire protocol.
type SubProtocol interface {
	Ready(sub *Subscriber) error
	Recv(r io.Reader) (*Message, error)
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

// checkSignature is called to check the protocol signature
// of a greeting message.
func checkSignature(sig [2]byte) error {
	if sig[0] != ProtocolSignature[0] || sig[1] != ProtocolSignature[1] {
		return errors.New("unexpected protocol signature")
	}
	return nil
}
