package pubsub

import (
	"io"
)

// subProtoV1 version 1.0 of the Sub protocol.
type subProtoV1 struct{}

// Ready is called to tell a Publisher that a Subscriber is ready
// to receive messages/data.
//
// TODO: support multiple topics(?)
func (s *subProtoV1) Ready(sub *Subscriber) error {
	greeting := greetingV1{}
	err := greeting.read(sub.conn)
	if err != nil {
		return err
	}

	ready := readyV1{
		Ready: Ready{
			Major: majorV1,
			Minor: minorV1,
		},
		props: make(map[string]string),
	}

	ready.props[propIDKey] = sub.id
	ready.props[propTopicKey] = sub.topics[0]

	return ready.write(sub.conn)
}

// Recv is called to receive the message from a Publisher.
func (s *subProtoV1) Recv(r io.Reader) (*Message, error) {
	msg := messageV1{}

	err := msg.readNull(r)
	if err != nil {
		return nil, err
	}

	if !msg.hasNull() {
		errMsg := errorV1{code: msg.null}
		err = errMsg.readAfterCode(r)
		if err != nil {
			return nil, err
		}

		return nil, errMsg
	}

	err = msg.readAfterNull(r)
	if err != nil {
		return nil, err
	}

	return &Message{
		Topic: string(msg.topic),
		Body:  msg.body,
	}, nil
}
