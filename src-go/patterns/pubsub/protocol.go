package pubsub

import (
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
)

// ProtocolSignature is used to identify messages belonging to the Pub-Sub protocol.
// 10101011 11[000001], where [nnnnnn] identifies the protocol.
var ProtocolSignature = [...]byte{0xAB, 0xC1}

const (
	// propIDKey is used as  of the 'id' message property.
	//propIDKey = "id"
	// propEndpointIDKey is the identifier of the remote endpoint.
	propEndpointIDKey = "epid"
	// propTopicKey is the identifier of the 'topic' message & queue properties.
	propTopicKey = "topic"
)

// PubProtocol is the Publisher-side of the pub-sub wire protocol.
type PubProtocol interface {
	core.GreetSendReceiver
}

// SubProtocol is the Subscriber-side of the pub-sub wire protocol.
type SubProtocol interface {
	core.GreetSendReceiver
}

// A Greeting is the first message sent by a Publisher to a Subscriber.
type Greeting struct {
	Signature [2]byte
	Major     uint8
	Minor     uint8
}

// A Ready message is sent by a Subscriber in response to a Greeting.
type Ready struct {
	Major uint8
	Minor uint8
}

// checkSignature is called to check the protocol signature
// of a greeting message.
func checkSignature(sig [2]byte) error {
	if sig != ProtocolSignature {
		return core.ErrInvalidSig
	}
	return nil
}
