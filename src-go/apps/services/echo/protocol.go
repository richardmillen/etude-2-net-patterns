package echo

import (
	"io"

	"github.com/richardmillen/etude-2-net-patterns/src-go/frames"
)

// protocolSignature is used to identify messages belonging to the Echo protocol.
var protocolSignature = [...]byte{0x00, 0x02}

// request echo message sent from client to server.
type request struct {
	signature  [2]byte
	endpointID [16]byte
	textLen    uint8
	text       []byte
}

func (req *request) read(r io.Reader) (err error) {
	req.signature, err = frames.ReadSig(r)
	if err != nil {
		return
	}

	epid, err := frames.ReadBytes(r, 16)
	if err != nil {
		return
	}
	copy(req.endpointID[:], epid)

	req.textLen, err = frames.ReadUInt8(r)
	if err != nil {
		return
	}

	req.text, err = frames.ReadBytes(r, req.textLen)
	return
}

func (req *request) write(w io.Writer) (err error) {
	buf := make([]byte, 2+16+1+req.textLen)
	bufView := buf

	bufView = frames.WriteBytes(bufView, req.signature[:])
	bufView = frames.WriteBytes(bufView, req.endpointID[:])
	bufView = frames.WriteUInt8(bufView, req.textLen)
	bufView = frames.WriteBytes(bufView, req.text)

	_, err = w.Write(buf)
	return
}

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

	rep.bodyLen, err := frames.ReadUInt8(r)
	if err != nil {
		return
	}

	rep.body, err = frames.ReadBytes(r, rep.bodyLen)
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
