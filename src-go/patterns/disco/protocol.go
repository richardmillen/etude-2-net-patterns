package disco

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/core"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

// minServiceNameLen is the minimum number of bytes of a service name.
const minServiceNameLen = 2

// protocolSignature is used to identify messages belonging to the discovery protocol.
// 10101011 11[000010], where [nnnnnn] identifies the protocol.
var protocolSignature = [...]byte{0xAB, 0xC2}

// ErrBadSurveyData ...
var ErrBadSurveyData = errors.New("invalid survey data")

const (
	cmdSurvey     uint8 = 0
	cmdResponseOK uint8 = 1
)

// checkSignature is called to check the protocol signature
// of a greeting message.
func checkSignature(sig [2]byte) error {
	if sig != protocolSignature {
		return core.ErrInvalidSig
	}
	return nil
}

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

func (s *survey) minDataLen() int {
	return 2
}

func (s *survey) maxDataLen() int {
	return 8
}

// toBytes turns a survey message into a slice of bytes.
// HACK: passing the min/max in as parameters is ugly!
func (s *survey) toBytes(minDataLen int, maxDataLen int) ([]byte, error) {
	if len(s.data) < minDataLen || len(s.data) > maxDataLen {
		return nil, ErrBadSurveyData
	}

	data := string(s.data)

	buf := make([]byte, 2+1+uuid.Size+len(data))
	bufView := buf

	copy(bufView, s.signature[:])
	bufView = bufView[2:]

	bufView[0] = byte(s.command)
	bufView = bufView[1:]

	copy(bufView, s.surveyID[:])
	bufView = bufView[16:]

	copy(bufView, []byte(data))

	return buf, nil
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

	s.surveyID = uuid.NewFrom(b[0:uuid.Size])
	b = b[uuid.Size:]

	s.data = make([]byte, len(b))
	copy(s.data, b)

	return
}

func (s *survey) readFrom(conn net.PacketConn) (addr net.Addr, err error) {
	log.Println("waiting for incoming survey...")

	buf := make([]byte, 2+1+uuid.Size+s.maxDataLen())
	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		return
	}

	if n < s.minDataLen() {
		err = ErrBadSurveyData
		return
	}

	s.fromBytes(buf[0:n])
	return
}

func (s *survey) write(w io.Writer) (err error) {
	var buf []byte

	buf, err = s.toBytes(s.minDataLen(), s.maxDataLen())
	if err != nil {
		return
	}

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

func (res *response) minDataLen() int {
	return 8
}

func (res *response) maxDataLen() int {
	return 64
}

// read is called to receive a survey response.
//
// this method validates the command and survey id because besides
// receiving a response to our survey message we could also receive:
// 	a) a copy of the original survey message (that we sent).
//	b) someone elses survey message.
// 	c) a response to someone elses survey message.
func (res *response) read(r io.Reader, surveyID uuid.Bytes) (err error) {
	buf := make([]byte, 2+1+uuid.Size+res.maxDataLen())
	n, err := r.Read(buf)
	if err != nil {
		return
	}

	if n < res.minDataLen() {
		err = ErrBadSurveyData
		return
	}

	res.fromBytes(buf[0:n])

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

func (res *response) writeTo(conn net.PacketConn, addr net.Addr) (err error) {
	log.Println("sending survey response...")
	buf, err := res.toBytes(res.minDataLen(), res.maxDataLen())
	_, err = conn.WriteTo(buf, addr)
	return
}
