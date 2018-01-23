package echo

import (
	"io"

	"github.com/richardmillen/etude-2-net-patterns/src-go/frames"
)

// protocolSignature is used to identify messages belonging to the Echo protocol.
// 10101011 11[000011], where [nnnnnn] identifies the protocol.
var protocolSignature = [...]byte{0xAB, 0xC3}

// request echo message sent from client to server.
type request struct {
	signature [2]byte
	textLen   uint8
	text      []byte
}

func (req *request) read(r io.Reader) (err error) {
	req.signature, err = frames.ReadSig(r)
	if err != nil {
		return
	}

	req.textLen, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	req.text, err = frames.ReadBytes(r, int64(req.textLen))
	return
}

func (req *request) write(w io.Writer) (err error) {
	buf := make([]byte, 2+1+req.textLen)
	bufView := buf

	bufView = frames.WriteBytes(bufView, req.signature[:])
	bufView = frames.WriteUInt8(bufView, req.textLen)
	bufView = frames.WriteBytes(bufView, req.text)

	_, err = w.Write(buf)
	return
}

const (
	codeOK  = 0
	codeErr = 1
)

// reply echo message sent from server back to client.
type reply struct {
	code    uint8
	bodyLen uint8
	body    []byte
}

func (rep *reply) read(r io.Reader) (err error) {
	rep.code, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	rep.bodyLen, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	rep.body, err = frames.ReadBytes(r, int64(rep.bodyLen))
	return
}

func (rep *reply) write(w io.Writer) (err error) {
	buf := make([]byte, 1+1+rep.bodyLen)
	bufView := buf

	bufView = frames.WriteUInt8(bufView, rep.code)
	bufView = frames.WriteUInt8(bufView, rep.bodyLen)
	bufView = frames.WriteBytes(bufView, rep.body)

	_, err = w.Write(buf)
	return
}
