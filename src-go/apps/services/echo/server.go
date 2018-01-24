package echo

import (
	"io"
	"log"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

// NewServer constructs a new instance of an Echo Server.
//
// The quit channel is buffered in order to avoid
// a deadlock when Server.Close() is called.
func NewServer(l net.Listener) *Server {
	s := &Server{
		UUID:     uuid.New(),
		listener: l,
		quit:     make(chan bool, 1),
		stopped:  make(chan bool),
	}

	go s.run()
	return s
}

// Server represents an instance of an Echo server.
type Server struct {
	UUID     uuid.Bytes
	listener net.Listener
	quit     chan bool
	stopped  chan bool
}

func (s *Server) run() {
	defer func() {
		log.Println("echo server stopped.")
		s.stopped <- true
	}()

	for {
		select {
		case <-s.quit:
			return
		default:
			log.Println("waiting for incoming echo request...")

			conn, err := s.listener.Accept()
			if check.Log(err) {
				continue
			}

			text, err := s.recv(conn)
			if check.Log(err) {
				continue
			}

			check.Log(s.send(conn, text))
		}
	}
}

// recv is called to receive an echo request message.
func (s *Server) recv(r io.Reader) (string, error) {
	req := request{}

	for {
		err := req.read(r)
		if err != nil {
			return "", err
		}
		return string(req.text), nil
	}
}

// send is called to send an echo reply message.
func (s *Server) send(w io.Writer, text string) error {
	rep := reply{
		code:    codeOK,
		bodyLen: uint8(len(text)),
		body:    []byte(text),
	}
	return rep.write(w)
}

// Addr returns the address of the echo server.
func (s *Server) Addr() string {
	return core.GetEndpointAddress(s.listener.Addr())
}

// Close quits the server background goroutine and closes the TCP connection.
//
// Server.quit is buffered so no need to use cumbersome select{}.
func (s *Server) Close() (err error) {
	s.quit <- true
	err = s.listener.Close()

	<-s.stopped
	return
}
