package pubsub

import (
	"fmt"
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

func (p *protoV1) send(topic string, data []byte) error {
	var msg message

}

func (p *protoV1) checkVersion(ready *ready) error {
	if ready.major == v1Major && ready.minor == v1Minor {
		return nil
	}
	return fmt.Errorf("protocol mismatch. version %d.%d required", v1Major, v1Minor)
}

// greeting is sent by the publisher immediately after a subscriber connects.
type greeting struct {
	signature [2]byte
	major     byte
	minor     byte
}

// sendTo sends the greeting to a subscription.
func (g *greeting) sendTo(sub *subscription) (err error) {
	msg := make([]byte, 4)
	buf := msg

	copy(buf, g.signature[:])
	buf = buf[len(g.signature):]

	buf[0] = g.major
	buf[1] = g.minor

	_, err = sub.conn.Write(msg)
	return
}

// ready is sent by the subscriber following a greeting.
type ready struct {
	major    byte
	minor    byte
	propsLen [2]byte
	props    []byte
}

// recvFrom receives a ready message from a subscription.
func (r *ready) recvFrom(sub *subscription) error {

}

func (r *ready) readProps() map[string]string {
	props := make(map[string]string)
	// TODO: read the props data in
	return props
}

// refuse is sent by the publisher when the subscriber failed to connect properly.
type refuse struct {
	code      byte
	reasonLen byte
	reason    []byte
}

// message is sent from publisher to topic subscriber.
type message struct {
	null     byte
	topicLen [2]byte
	topic    []byte
	dataLen  [2]byte
	data     []byte
}
