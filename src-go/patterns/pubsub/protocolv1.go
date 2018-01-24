package pubsub

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/richardmillen/etude-2-net-patterns/src-go/frames"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

const (
	// majorV1 protocol major version number.
	majorV1 = 1
	// minorV1 protocol minor version number.
	minorV1 = 0
)

// checkVersion is called to check a version number is ok.
func checkVersion(major uint8, minor uint8) error {
	if major != majorV1 || minor != minorV1 {
		return fmt.Errorf("protocol mismatch. version %d.%d required", majorV1, minorV1)
	}
	return nil
}

// pubProtoV1 is version 1.0 of the Pub protocol.
type pubProtoV1 struct{}

// Greet is called by a Publisher to set up a new connection with a Subscriber.
func (p *pubProtoV1) Greet(q *core.Queue) error {
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
	q.SetProp(core.PropUUIDKey, uuid.NewFrom(ready.props[core.PropUUIDKey]))
	q.SetProp(propTopicKey, string(ready.props[propTopicKey]))

	return nil
}

// Send is called to send content to a Subscriber.
// Only perform the send if the queue is subscribing to the topic.
func (p *pubProtoV1) Send(q *core.Queue, v interface{}) error {
	msg := v.(*Message)

	if !p.isQueueSubscribing(q, msg.Topic) {
		return nil
	}

	var m messageV1
	m.topicLen = uint16(len(msg.Topic))
	m.topic = []byte(msg.Topic)
	m.bodyLen = uint16(len(msg.Body))
	m.body = msg.Body

	return m.write(q.Conn())
}

// Recv is not supported.
func (p *pubProtoV1) Recv(q *core.Queue) (interface{}, error) {
	return nil, errors.New("pubProtoV1.Recv: not supported")
}

func (p *pubProtoV1) isQueueSubscribing(q *core.Queue, topic string) bool {
	return strings.HasPrefix(q.Prop(propTopicKey).(string), topic)
}

// subProtoV1 is version 1.0 of the Sub protocol.
type subProtoV1 struct{}

// Greet is called by a Subscriber to respond to a Publisher's Greeting,
// informing the Publisher that it's ready to receive topic messages/data.
//
// TODO: support multiple topics(?)
func (s *subProtoV1) Greet(q *core.Queue) error {
	greeting := greetingV1{}
	err := greeting.read(q.Conn())
	if err != nil {
		return err
	}

	ready := readyV1{
		Ready: Ready{
			Major: majorV1,
			Minor: minorV1,
		},
		props: make(map[string][]byte),
	}

	ready.props[core.PropUUIDKey] = q.Prop(core.PropUUIDKey).(uuid.Bytes)
	ready.props[propTopicKey] = []byte(q.Prop(propTopicKey).(string))

	return ready.write(q.Conn())
}

// Send is not supported.
func (s *subProtoV1) Send(q *core.Queue, v interface{}) error {
	return errors.New("subProtoV1.Send: not supported")
}

// Recv is called to receive the message from a Publisher.
func (s *subProtoV1) Recv(q *core.Queue) (interface{}, error) {
	msg := messageV1{}

	err := msg.readNull(q.Conn())
	if err != nil {
		return nil, err
	}

	if !msg.hasNull() {
		errMsg := errorV1{code: msg.null}
		err = errMsg.readAfterCode(q.Conn())
		if err != nil {
			return nil, err
		}

		return nil, errMsg
	}

	err = msg.readAfterNull(q.Conn())
	if err != nil {
		return nil, err
	}

	return &Message{
		Topic: string(msg.topic),
		Body:  msg.body,
	}, nil
}

// A greetingV1 message is the first message sent by a v1.0 publisher to a subscriber.
type greetingV1 struct {
	Greeting
}

func (msg *greetingV1) read(r io.Reader) (err error) {
	msg.Signature, err = frames.ReadSig(r)
	if err != nil {
		return
	}

	err = checkSignature(msg.Signature)
	if err != nil {
		return
	}

	msg.Major, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	msg.Minor, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	err = checkVersion(msg.Major, msg.Minor)
	if err != nil {
		return
	}
	return
}

// write sends the greeting to a subscription.
func (msg *greetingV1) write(w io.Writer) (err error) {
	buf := make([]byte, 4)

	copy(buf, msg.Signature[:])
	buf[2] = byte(msg.Major)
	buf[3] = byte(msg.Minor)

	_, err = w.Write(buf)
	return
}

// A readyV1 message is sent by a v1.0 subscriber in response to a greeting.
type readyV1 struct {
	Ready
	propsLen uint16
	props    map[string][]byte
}

// read receives a ready message from a subscription.
func (msg *readyV1) read(r io.Reader) (err error) {
	msg.Major, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	msg.Minor, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	err = checkVersion(msg.Major, msg.Minor)
	if err != nil {
		return
	}

	msg.propsLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	log.Println("ready msg props len:", msg.propsLen)

	msg.props, err = frames.ReadProps(r, int64(msg.propsLen))
	return
}

func (msg *readyV1) write(w io.Writer) (err error) {
	props := frames.PropsToBytes(msg.props)
	msg.propsLen = uint16(len(props))

	log.Println("ready msg props len:", msg.propsLen)

	buf := make([]byte, 1+1+2+len(props))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, msg.Major)
	bufView = frames.WriteUInt8(bufView, msg.Minor)
	bufView = frames.WriteUInt16(bufView, msg.propsLen)
	bufView = frames.WriteBytes(bufView, props)

	_, err = w.Write(buf)
	return
}

// An errorV1 message is sent by a v1.0 publisher when the subscriber failed to connect properly.
type errorV1 struct {
	code      uint8
	reasonLen uint16
	reason    []byte
}

func (msg errorV1) Error() string {
	return fmt.Sprintf("%s (%d)", msg.reason, msg.code)
}

func (msg *errorV1) readAfterCode(r io.Reader) (err error) {
	msg.reasonLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	msg.reason, err = frames.ReadBytes(r, int64(msg.reasonLen))
	return
}

func (msg *errorV1) write(w io.Writer) (err error) {
	buf := make([]byte, 1+2+len(msg.reason))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, msg.code)
	bufView = frames.WriteUInt16(bufView, msg.reasonLen)
	bufView = frames.WriteBytes(bufView, msg.reason)

	_, err = w.Write(buf)
	return
}

// A messageV1 is sent from a v1.0 Publisher to topic Subscriber.
type messageV1 struct {
	// null is used to differentiate a message from a 'refuse'.
	null     uint8
	topicLen uint16
	topic    []byte
	bodyLen  uint16
	body     []byte
}

func (msg *messageV1) hasNull() bool {
	return msg.null == 0
}

// readNull reads the first byte (uint8) which is expected to be
// null. if it's not te the Publisher is reporting an error.
func (msg *messageV1) readNull(r io.Reader) (err error) {
	msg.null, err = frames.ReadUInt8(r)
	return
}

func (msg *messageV1) readAfterNull(r io.Reader) (err error) {
	msg.topicLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	msg.topic, err = frames.ReadBytes(r, int64(msg.topicLen))
	if err != nil {
		return
	}

	msg.bodyLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	msg.body, err = frames.ReadBytes(r, int64(msg.bodyLen))
	if err != nil {
		return
	}

	return
}

// write is called to send a message to a subscription endpoint (Subscriber).
func (msg *messageV1) write(w io.Writer) (err error) {
	buf := make([]byte, 1+2+len(msg.topic)+2+len(msg.body))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, msg.null)
	bufView = frames.WriteUInt16(bufView, msg.topicLen)
	bufView = frames.WriteBytes(bufView, msg.topic)
	bufView = frames.WriteUInt16(bufView, msg.bodyLen)
	bufView = frames.WriteBytes(bufView, msg.body)

	_, err = w.Write(buf)
	return
}
