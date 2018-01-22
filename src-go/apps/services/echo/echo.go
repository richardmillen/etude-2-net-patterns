package echo

import (
	"fmt"
	"io"

	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/disco"
)

// ServiceName is the name used when referring to the Echo service.
const ServiceName = "echo"

// Send is called by a client to send an echo request to an echo server.
func Send(w io.Writer, endpoint *disco.Endpoint, text string) error {
	req := request{}
	req.signature = protocolSignature
	req.endpointID = endpoint.UUID
	req.textLen = len(text)
	req.text = text
	return req.write(w)
}

// Recv is called by a client to receive an echo response from an echo server.
func Recv(r io.Reader) (string, error) {
	rep := reply{}

	err := rep.read(r)
	if err != nil {
		return "", err
	}

	if rep.code != 0 {
		return "", fmt.Errorf("%s (%d)", rep.body, rep.code)
	}

	return string(rep.body), nil
}
