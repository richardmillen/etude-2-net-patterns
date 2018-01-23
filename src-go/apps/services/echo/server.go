package echo

import (
	"io"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

// NewServer constructs a new instance of an Echo Server.
func NewServer(l net.Listener) *Server {
	s := &Server{
		UUID: uuid.New(),
		quit: make(chan bool),
	}
	go s.run()
	return s
}

// Server represents an instance of an Echo server.
type Server struct {
	UUID     uuid.Bytes
	listener net.Listener
	quit     chan bool
}

func (s *Server) run() {
	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			check.Error(err)

			text, err := s.recv(conn)
			check.Error(err)

			check.Must(s.send(conn, text))
		}
	}
}

// Addr returns the address of the echo server.
func (s *Server) Addr() string {
	return s.listener.Addr().String()
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

// Close quits the server background goroutine and closes the TCP connection.
func (s *Server) Close() error {
	s.quit <- true
	return s.listener.Close()
}
