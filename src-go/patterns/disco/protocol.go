package disco

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"

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
// 10101011 11[000010], where [nnnnnn] identifies the protocol.
var protocolSignature = [...]byte{0xAB, 0xC2}

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
	surveyID  uuid.Bytes
	data      []byte
}

// toBytes turns a survey message into a slice of bytes.
func (s *survey) toBytes() []byte {
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

	return buf
}

// fromBytes fills the fields of a survey message from a slice of bytes.
func (s *survey) fromBytes(b []byte) (err error) {
	copy(s.signature[:], b[0:2])
	b = b[2:]

	err = checkSignature(s.signature)
	if err != nil {
		return
	}

	s.command = uint8(b[0])
	b = b[1:]

	copy(s.surveyID[:], b[0:16])
	b = b[16:]

	s.data = make([]byte, len(b))
	copy(s.data, b)

	return
}

func (s *survey) readFrom(conn net.PacketConn) (addr net.Addr, err error) {
	buf := make([]byte, 2+1+16+maxServiceNameLen)
	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		return
	}
	check.IsGreaterEqual(n, minSurveyMsgLen, "survey message length")
	s.fromBytes(buf[0:n])
	return
}

func (s *survey) read(r io.Reader) (err error) {
	buf := make([]byte, 2+1+16+maxServiceNameLen)
	n, err := r.Read(buf)
	if err != nil {
		return
	}
	check.IsGreaterEqual(n, minSurveyMsgLen, "survey message length")
	s.fromBytes(buf[0:n])
	return
}

// TODO: refactor duplicate code (write / writeTo)
func (s *survey) writeTo(conn net.PacketConn, addr net.Addr) (err error) {
	buf := s.toBytes()
	_, err = conn.WriteTo(buf, addr)
	return
}

// TODO: refactor duplicate code (write / writeTo)
func (s *survey) write(w io.Writer) (err error) {
	buf := s.toBytes()
	_, err = w.Write(buf)
	return
}

// ErrUnkSurveyID is returned when a survey response is received
// with a survey id that doesn't match the original survey request.
var ErrUnkSurveyID = errors.New("response.read: unknown survey id")

func newResponse(req *survey, addr string) *response {
	return &response{
		survey: survey{
			signature: req.signature,
			command:   cmdResponseOK,
			surveyID:  req.surveyID,
			data:      []byte(addr),
		},
	}
}

type response struct {
	survey
}

// read is called to receive a survey response.
//
// this method validates the command and survey id because besides
// receiving a response to our survey message we could also receive:
// 	a) a copy of the original survey message (that we sent).
//	b) someone elses survey message.
// 	c) a response to someone elses survey message.
func (res *response) read(r io.Reader, surveyID uuid.Bytes) (err error) {
	err = res.survey.read(r)
	if err != nil {
		return
	}

	if res.command != cmdResponseOK {
		err = fmt.Errorf("%s (%d)", string(res.data), res.command)
		return
	}

	if !uuid.Equal(res.surveyID, surveyID) {
		err = ErrUnkSurveyID
		return
	}

	return
}

// readFrom is called to receive a survey response.
//
// this method validates the command and survey id because besides
// receiving a response to our survey message we could also receive:
// 	a) a copy of the original survey message (that we sent).
//	b) someone elses survey message.
// 	c) a response to someone elses survey message.
/*func (res *response) readFrom(conn net.PacketConn, surveyID [16]byte) (addr net.Addr, err error) {
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
}*/

// writeTo
/*func (res *response) writeTo(conn net.PacketConn, addr net.Addr) (err error) {
}*/
