package pubsub

import (
	"fmt"

	"github.com/richardmillen/etude-2-net-patterns/src-go/frames"
)

var v1Signature = [...]byte{0x01, 0x01}

const (
	v1Major = 1
	v1Minor = 0
)

const (
	v1PropID    = "id"
	v1PropTopic = "topic"
)

// protoV1 is pub-sub protocol version 1.0.
type protoV1 struct{}

func (p *protoV1) greet(sub *subscription) error {
	greeting := greeting{
		signature: v1Signature,
		major:     v1Major,
		minor:     v1Minor,
	}

	err := greeting.sendTo(sub)
	if err != nil {
		return err
	}

	ready := ready{}
	err = ready.recvFrom(sub)
	if err != nil {
		return err
	}

	err = p.checkVersion(&ready)
	if err != nil {
		return err
	}

	sub.proto = p

	props := ready.readProps()
	sub.id = props[v1PropID]
	sub.topic = props[v1PropTopic]

	return nil
}

func (p *protoV1) sendTo(sub *subscription, topic string, body []byte) (err error) {
	var msg message

	msg.topicLen = uint16(len(topic))
	msg.topic = []byte(topic)
	msg.bodyLen = uint16(len(body))
	msg.body = body

	return msg.sendTo(sub)
}

func (p *protoV1) checkVersion(ready *ready) error {
	if ready.major == v1Major && ready.minor == v1Minor {
		return nil
	}
	return fmt.Errorf("protocol mismatch. version %d.%d required", v1Major, v1Minor)
}

// greeting is a message sent by the publisher immediately after a subscriber connects.
type greeting struct {
	signature [2]byte
	major     uint8
	minor     uint8
}

// sendTo sends the greeting to a subscription.
func (g *greeting) sendTo(sub *subscription) (err error) {
	buf := make([]byte, 4)

	copy(buf, g.signature[:])
	buf[2] = byte(g.major)
	buf[3] = byte(g.minor)

	_, err = sub.conn.Write(buf)
	return
}

// ready is a message sent by the subscriber in response to a greeting.
type ready struct {
	major    uint8
	minor    uint8
	propsLen uint16
	props    []byte
}

// recvFrom receives a ready message from a subscription.
func (r *ready) recvFrom(sub *subscription) (err error) {
	r.major, err = frames.ReadUInt8(sub.conn)
	if err != nil {
		return
	}

	r.minor, err = frames.ReadUInt8(sub.conn)
	if err != nil {
		return
	}

	r.propsLen, err = frames.ReadUInt16(sub.conn)
	if err != nil {
		return
	}

	r.props, err = frames.ReadBytes(sub.conn, int64(r.propsLen))
	return
}

func (r *ready) readProps() map[string]string {
	props := make(map[string]string)
	// TODO: read the props data in
	return props
}

// refuse is a message sent by the publisher when the subscriber failed to connect properly.
type refuse struct {
	// code contains an identifier for the error.
	code      int8
	reasonLen uint8
	reason    []byte
}

// message is sent from publisher to topic subscriber.
type message struct {
	// null is used to differentiate a message from a 'refuse'.
	null     uint8
	topicLen uint16
	topic    []byte
	bodyLen  uint16
	body     []byte
}

// sendTo is called to send a message to a subscription endpoint (Subscriber).
func (m *message) sendTo(sub *subscription) (err error) {
	buf := make([]byte, 1+2+len(m.topic)+2+len(m.body))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, m.null)
	bufView = frames.WriteUInt16(bufView, m.topicLen)
	bufView = frames.WriteBytes(bufView, m.topic)
	bufView = frames.WriteUInt16(bufView, m.bodyLen)
	bufView = frames.WriteBytes(bufView, m.body)

	_, err = sub.conn.Write(buf)
	return
}
