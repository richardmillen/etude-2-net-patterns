package pubsub

const (
	// majorV1 protocol major version number.
	majorV1 = 1
	// minorV1 protocol minor version number.
	minorV1 = 0
)

type greetingV1 struct {
	Greeting
}

// readyV1 is a message sent by the subscriber in response to a greeting.
type readyV1 struct {
	Ready
	propsLen uint16
	props    []byte
}

// refuseV1 is a message sent by the publisher when the subscriber failed to connect properly.
type refuseV1 struct {
	// code contains an identifier for the error.
	code      int8
	reasonLen uint8
	reason    []byte
}

// messageV1 is sent from publisher to topic subscriber.
type messageV1 struct {
	// null is used to differentiate a message from a 'refuse'.
	null     uint8
	topicLen uint16
	topic    []byte
	bodyLen  uint16
	body     []byte
}
