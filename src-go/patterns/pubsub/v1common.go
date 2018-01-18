package pubsub

import (
	"fmt"
	"io"

	"github.com/richardmillen/etude-2-net-patterns/src-go/frames"
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

type greetingV1 struct {
	Greeting
}

func (v1 *greetingV1) read(r io.Reader) (err error) {
	v1.Signature, err = frames.ReadSig(r)
	if err != nil {
		return
	}

	err = checkSignature(v1.Signature)
	if err != nil {
		return
	}

	v1.Major, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	v1.Minor, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	err = checkVersion(v1.Major, v1.Minor)
	if err != nil {
		return
	}
	return
}

// write sends the greeting to a subscription.
func (v1 *greetingV1) write(conn io.Writer) (err error) {
	buf := make([]byte, 4)

	copy(buf, v1.Signature[:])
	buf[2] = byte(v1.Major)
	buf[3] = byte(v1.Minor)

	_, err = conn.Write(buf)
	return
}

// readyV1 is a message sent by the subscriber in response to a greeting.
type readyV1 struct {
	Ready
	propsLen uint16
	//props    []byte
	props map[string]string
}

// read receives a ready message from a subscription.
func (v1 *readyV1) read(r io.Reader) (err error) {
	v1.Major, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	v1.Minor, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	err = checkVersion(v1.Major, v1.Minor)
	if err != nil {
		return
	}

	v1.propsLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	v1.props, err = frames.ReadProps(r, int64(v1.propsLen))
	return
}

func (v1 *readyV1) write(w io.Writer) (err error) {
	props := frames.PropsToBytes(v1.props)
	v1.propsLen = uint16(len(props))

	buf := make([]byte, 1+1+2+len(props))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, v1.Major)
	bufView = frames.WriteUInt8(bufView, v1.Minor)
	bufView = frames.WriteUInt16(bufView, v1.propsLen)
	bufView = frames.WriteBytes(bufView, props)

	_, err = w.Write(buf)
	return
}

// errorV1 is a message sent by the publisher when the subscriber failed to connect properly.
type errorV1 struct {
	// code contains an identifier for the error.
	code      uint8
	reasonLen uint16
	reason    []byte
}

func (v1 errorV1) Error() string {
	return fmt.Sprintf("%s (%d)", v1.reason, v1.code)
}

func (v1 *errorV1) readAfterCode(r io.Reader) (err error) {
	v1.reasonLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	v1.reason, err = frames.ReadBytes(r, int64(v1.reasonLen))
	return
}

func (v1 *errorV1) write(w io.Writer) (err error) {
	buf := make([]byte, 1+2+len(v1.reason))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, v1.code)
	bufView = frames.WriteUInt16(bufView, v1.reasonLen)
	bufView = frames.WriteBytes(bufView, v1.reason)

	_, err = w.Write(buf)
	return
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

func (v1 *messageV1) hasNull() bool {
	return v1.null == 0
}

// readNull reads the first byte (uint8) which is expected to be
// null. if it's not te the Publisher is reporting an error.
func (v1 *messageV1) readNull(r io.Reader) (err error) {
	v1.null, err = frames.ReadUInt8(r)
	return
}

func (v1 *messageV1) readAfterNull(r io.Reader) (err error) {
	v1.topicLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	v1.topic, err = frames.ReadBytes(r, int64(v1.topicLen))
	if err != nil {
		return
	}

	v1.bodyLen, err = frames.ReadUInt16(r)
	if err != nil {
		return
	}

	v1.body, err = frames.ReadBytes(r, int64(v1.bodyLen))
	if err != nil {
		return
	}

	return
}

// write is called to send a message to a subscription endpoint (Subscriber).
func (v1 *messageV1) write(w io.Writer) (err error) {
	buf := make([]byte, 1+2+len(v1.topic)+2+len(v1.body))
	bufView := buf

	bufView = frames.WriteUInt8(bufView, v1.null)
	bufView = frames.WriteUInt16(bufView, v1.topicLen)
	bufView = frames.WriteBytes(bufView, v1.topic)
	bufView = frames.WriteUInt16(bufView, v1.bodyLen)
	bufView = frames.WriteBytes(bufView, v1.body)

	_, err = w.Write(buf)
	return
}
