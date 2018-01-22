package disco

import (
	"errors"
	"fmt"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns"
)

// minServiceNameLen is the minimum number of bytes of a service name.
const minServiceNameLen = 2

// maxServiceNameLen is the maximum number of bytes of a service name.
const maxServiceNameLen = 8

// maxEndpointDataLen is the maximun number of bytes of the endpoint info within a response message.
const maxEndpointDataLen = 64

// minSurveyMsgLen represents the minimum valid length of a survey request message.
// n.b. shown by adding the lengths of a survey messages component parts.
const minSurveyMsgLen = 2 + 1 + 16 + minServiceNameLen

// maxSurveyMsgLen represents the maximum valid length of a survey request message.
const maxSurveyMsgLen = 2 + 1 + 16 + maxServiceNameLen

// protocolSignature is used to identify messages belonging to the discovery protocol.
var protocolSignature = [...]byte{0x02, 0x00}

// checkSignature is called to check the protocol signature
// of a greeting message.
func checkSignature(sig [2]byte) error {
	if sig != protocolSignature {
		return patterns.ErrInvalidSig
	}
	return nil
}

const (
	cmdSurvey     uint8 = 0
	cmdResponseOK uint8 = 1
)

// survey is both a message sent by a Surveyor when searching for a service
// and also a response message sent by a Candidate.
//
// the 'data' field is 'maxServiceNameLen' in a survey message, but is
// 'maxEndpointDataLen' in a response message.
type survey struct {
	signature [2]byte
	command   uint8
	surveyID  [16]byte
	data      []byte
}

func (s *survey) readFrom(conn net.PacketConn) (addr net.Addr, err error) {
	buf := make([]byte, 2+1+16+maxServiceNameLen)
	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		return
	}
	check.IsGreaterEqual(n, minSurveyMsgLen, "survey message length")

	bufView := buf[0:n]

	copy(s.signature[:], bufView[0:2])
	bufView = bufView[2:]

	err = checkSignature(s.signature)
	if err != nil {
		return
	}

	s.command = uint8(bufView[0])
	bufView = bufView[1:]

	copy(s.surveyID[:], bufView[0:16])
	bufView = bufView[16:]

	s.data = make([]byte, len(bufView))
	copy(s.data, bufView)

	return
}

func (s *survey) writeTo(conn net.PacketConn, addr net.Addr) (err error) {
	serviceName := string(s.data[:])
	check.IsInRange(len(serviceName), minServiceNameLen, maxServiceNameLen, "service name length")

	buf := make([]byte, 2+1+16+len(serviceName))
	bufView := buf

	copy(bufView, s.signature[:])
	bufView = bufView[2:]

	bufView[0] = byte(s.command)
	bufView = bufView[1:]

	copy(bufView, s.surveyID[:])
	bufView = bufView[16:]

	copy(bufView, []byte(serviceName))

	_, err = conn.WriteTo(buf, addr)
	return
}

// ErrUnkSurveyID is returned when a survey response is received
// with a survey id that doesn't match the original survey request.
var ErrUnkSurveyID = errors.New("response.read: unknown survey id")

func newResponse(req *survey) *response {

}

type response struct {
	survey
}

func (res *response) setEndpoint(endpoint *Endpoint) {

}

// readFrom is called to receive a survey response.
//
// this method validates the command and survey id because besides
// receiving a response to our survey message we could also receive:
// 	a) a copy of the original survey message (that we sent).
//	b) someone elses survey message.
// 	c) a response to someone elses survey message.
func (res *response) readFrom(conn net.PacketConn, surveyID [16]byte) (addr net.Addr, err error) {
	addr, err = res.survey.readFrom(conn)
	if err != nil {
		return
	}

	if res.command != cmdResponseOK {
		err = fmt.Errorf("%s (%d)", string(res.data), res.command)
		return
	}

	if res.surveyID != surveyID {
		err = ErrUnkSurveyID
		return
	}

	return
}

// writeTo
func (res *response) writeTo(conn net.PacketConn, addr net.Addr) (err error) {
}
