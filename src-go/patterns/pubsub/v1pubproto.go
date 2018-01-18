package pubsub

import (
	"io"
)

// pubProtoV1 version 1.0 of the Pub protocol.
type pubProtoV1 struct{}

// Greet is called by a Publisher to set up a new connection with a Subscriber.
func (p *pubProtoV1) Greet(q *Queue) error {
	greeting := greetingV1{
		Greeting{
			Signature: ProtocolSignature,
			Major:     majorV1,
			Minor:     minorV1,
		},
	}

	err := greeting.write(q.Conn())
	if err != nil {
		return err
	}

	ready := readyV1{}
	err = ready.read(q.Conn())
	if err != nil {
		return err
	}

	q.SetProtocol(p)
	q.SetProps(ready.props)

	return nil
}

// Send is called to send content to a Subscriber.
func (p *pubProtoV1) Send(conn io.Writer, m *Message) error {
	var msg messageV1

	msg.topicLen = uint16(len(m.Topic))
	msg.topic = []byte(m.Topic)
	msg.bodyLen = uint16(len(m.Body))
	msg.body = m.Body

	return msg.write(conn)
}
