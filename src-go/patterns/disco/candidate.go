package disco

import (
	"errors"
	"net"
	"strconv"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// NewCandidate constructs a new survey candidate.
func NewCandidate() Candidate {
	return Candidate{
		Port:             surveyPort,
		serviceEndpoints: make(map[string]*Endpoint),
		quit:             make(chan bool),
	}
}

// Candidate is a participant in a survey (service discovery).
// TODO: come up with better name (and description).
type Candidate struct {
	Port             int
	serviceEndpoints map[string]*Endpoint
	conn             net.PacketConn
	isOpen           bool
	quit             chan bool
}

// AddService is called to add a service name / endpoint address mapping.
func (c *Candidate) AddService(name string, addr string) {
	check.IsFalse(c.isOpen, "Candidate.isOpen")

	c.serviceEndpoints[name] = NewEndpoint(addr)
}

// Open is called to start responding to survey requests.
func (c *Candidate) Open() (err error) {
	if c.isOpen {
		return errors.New("candidate is already open for surveys")
	}

	c.conn, err = net.ListenPacket("udp", ":"+strconv.Itoa(c.Port))
	if check.Log(err) {
		return
	}

	c.isOpen = true
	go func() {
		for {
			select {
			case <-c.quit:
				c.conn.Close()
				return
			default:
				req := survey{}

				addr, err := req.readFrom(c.conn)
				if check.Log(err) {
					continue
				}

				// TODO: move this test into survey.readFrom():
				if req.command != cmdSurvey {
					continue
				}

				endpoint, ok := c.serviceEndpoints[string(req.data)]
				if !ok {
					continue
				}

				res := newResponse(&req)
				res.setEndpoint(endpoint)

				err = res.writeTo(c.conn, addr)
				if check.Log(err) {
					continue
				}
			}
		}
	}()

	return
}

// Close causes the survey candidate to stop listening for incoming survey requests.
// TODO: handle multiple calls(?)
func (c *Candidate) Close() error {
	c.quit <- true
	return c.conn.Close()
}
