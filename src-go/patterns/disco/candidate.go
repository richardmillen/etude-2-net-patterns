package disco

import (
	"log"
	"net"
	"strconv"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

// NewCandidate constructs a new survey candidate.
//
// The quit channel is buffered in order to avoid
// a deadlock when Candidate.Close() is called.
func NewCandidate() *Candidate {
	return &Candidate{
		Port:             surveyPort,
		serviceEndpoints: make(map[string]*Endpoint),
		quit:             make(chan bool, 1),
		stopped:          make(chan bool),
	}
}

// Candidate is a participant in a survey (service discovery).
// TODO: come up with better name (and description).
type Candidate struct {
	Port             int
	serviceEndpoints map[string]*Endpoint
	conn             net.PacketConn
	quit             chan bool
	stopped          chan bool
}

// AddService is called to add a service name / endpoint address mapping.
func (c *Candidate) AddService(name string, addr string) {
	c.serviceEndpoints[name] = NewEndpoint(addr)
}

// Open is called to start responding to survey requests.
func (c *Candidate) Open() (err error) {
	c.conn, err = net.ListenPacket("udp", ":"+strconv.Itoa(c.Port))
	if check.Log(err) {
		return
	}

	go func() {
		defer func() {
			log.Println("service candidate stopped.")
			c.stopped <- true
		}()

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

				res := newResponse(&req, endpoint.Addr)
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
//
// Candidate.quit is buffered so no need to use cumbersome select{}.
func (c *Candidate) Close() (err error) {
	c.quit <- true
	err = c.conn.Close()

	<-c.stopped
	return
}
