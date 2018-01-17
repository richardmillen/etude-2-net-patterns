package pubsub

import (
	"fmt"
	"io"

	"github.com/richardmillen/etude-2-net-patterns/src-go/frames"
)

// pubProtoV1 version 1.0 of the Pub protocol.
type pubProtoV1 struct{}

// Greet is called by a Publisher to set up a new connection with a Subscriber.
func (p *pubProtoV1) Greet(q SubscriptionQueue) error {
	greeting := greetingV1{
		Greeting{
			Signature: ProtocolSignature,
			Major:     majorV1,
			Minor:     minorV1,
		},
	}

	err := greeting.sendTo(q.Conn())
	if err != nil {
		return err
	}

	ready := readyV1{}
	err = ready.recvFrom(q.Conn())
	if err != nil {
		return err
	}

	err = p.checkVersion(&ready)
	if err != nil {
		return err
	}

	q.SetProtocol(p)
	q.SetProps(ready.readProps())

	return nil
}

// Send is called to send content to a Subscriber.
func (p *pubProtoV1) Send(conn io.ReadWriter, topic string, body []byte) error {
	var msg messageV1

	msg.topicLen = uint16(len(topic))
	msg.topic = []byte(topic)
	msg.bodyLen = uint16(len(body))
	msg.body = body

	return msg.sendTo(conn)
}

func (p *pubProtoV1) checkVersion(ready *readyV1) error {
	if ready.Major == majorV1 && ready.Minor == minorV1 {
		return nil
	}
	return fmt.Errorf("protocol mismatch. version %d.%d required", majorV1, minorV1)
}

// sendTo sends the greeting to a subscription.
func (g *greetingV1) sendTo(conn io.ReadWriter) (err error) {
	buf := make([]byte, 4)

	copy(buf, g.Signature[:])
	buf[2] = byte(g.Major)
	buf[3] = byte(g.Minor)

	_, err = conn.Write(buf)
	return
}

// recvFrom receives a ready message from a subscription.
func (r *readyV1) recvFrom(conn io.ReadWriter) (err error) {
	r.Major, err = frames.ReadUInt8(conn)
	if err != nil {
		return
	}

	r.Minor, err = frames.ReadUInt8(conn)
	if err != nil {
		return
	}

	r.propsLen, err = frames.ReadUInt16(conn)
	if err != nil {
		return
	}

	r.props, err = frames.ReadBytes(conn, int64(r.propsLen))
	return
}

func (r *readyV1) readProps() map[string]string {
	props := make(map[string]string)
	// TODO: read the props data in
	return props
}

// sendTo is called to send a message to a subscription endpoint (Subscriber).
func (m *messageV1) sendTo(conn io.ReadWriter) (err error) {
	buf := make([]byte, 1+2+len(m.topic)+2+len(m.body))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, m.null)
	bufView = frames.WriteUInt16(bufView, m.topicLen)
	bufView = frames.WriteBytes(bufView, m.topic)
	bufView = frames.WriteUInt16(bufView, m.bodyLen)
	bufView = frames.WriteBytes(bufView, m.body)

	_, err = conn.Write(buf)
	return
}
