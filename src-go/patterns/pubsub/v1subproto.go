package pubsub

import "io"

// subProtoV1 version 1.0 of the Sub protocol.
type subProtoV1 struct{}

// Ready is called to tell a Publisher that a Subscriber is ready
// to receive messages/data.
func (s *subProtoV1) Ready(sub *Subscriber) error {

}

// Recv is called to receive the message from a Publisher.
func (s *subProtoV1) Recv(conn io.ReadWriter) (topic string, body []byte, err error) {

}
