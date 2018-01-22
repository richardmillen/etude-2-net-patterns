package echo

import (
	"net"
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
type Server type {
	UUID string
	listener net.Listener
	quit chan bool
}

func (s *Server) run() {
	for {
		select {
		case <-quit:
			return
		default:
			conn, err := s.listener.Listen()
			check.Error(err)

			text, err := s.recv(conn)
			check.Error(err)

			check.Must(s.send(conn, text))
		}
	}
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

// recv is called to receive an echo request message.
func (s *Server) recv(r io.Reader) (string, error) {
	req := request{}

	for {
		check.Must(req.read(r))
		if string(req.endpointID[:]) == s.UUID {
			return string(req.text)
		}
	}
}

// send is called to send an echo reply message.
func (s *Server) send(w io.Writer, text string) error {
	rep := reply{}
	rep.bodyLen = len(text)
	rep.body = []byte(text)
	return rep.write(w)
}

func (s *Server) Close() error {
	s.quit <- true
	s.listener.Close()
}