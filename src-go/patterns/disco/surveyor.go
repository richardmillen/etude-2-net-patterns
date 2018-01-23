package disco

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

const surveyPort = 5677

// SurveyResponseFunc is called each time a survey respondant's message is received.
type SurveyResponseFunc = func(endpoint *Endpoint) error

// ErrEndSurvey is a special sentinal error value returned by a SurveyResponseFunc
// which is intended to tell a Surveyor to end an ongoing survey.
var ErrEndSurvey = errors.New("disco.Survey: end survey")

// NewSurveyor constructs a new Surveyor.
func NewSurveyor(addr string) *Surveyor {
	return &Surveyor{
		addr: addr,
		port: surveyPort,
		quit: make(chan bool),
	}
}

// Surveyor surveys the network for a service.
type Surveyor struct {
	addr string
	port int
	conn net.Conn
	quit chan bool
}

// Addr returns the full address to be surveyed.
func (s *Surveyor) Addr() string {
	return fmt.Sprintf("%s:%d", s.addr, s.port)
}

// Survey looks for a service by name, calling responseFunc for every response received within a specified timeframe.
// TODO: trap / enable(?) multiple calls to Survey().
func (s *Surveyor) Survey(responseFunc SurveyResponseFunc, timeout time.Duration, service string) (err error) {
	s.conn, err = net.Dial("udp", s.Addr())
	if err != nil {
		return
	}

	go func() {
		defer s.conn.Close()

		req := survey{}
		req.signature = protocolSignature
		req.command = cmdSurvey
		req.surveyID = uuid.New()
		req.data = []byte(service)

		err := req.write(s.conn)
		if check.Log(err) {
			return
		}

		timer := time.NewTimer(timeout)
		for {
			select {
			case <-timer.C:
				err = s.conn.Close()
				check.Log(err)
				return
			default:
				res := response{}
				err = res.read(s.conn, req.surveyID)
				if check.Log(err) {
					continue
				}

				// TODO: get endpoint info from response.
			}
		}
	}()

	return
}

// Close ends any ongoing surveys.
func (s *Surveyor) Close() error {
	close(s.quit)
	return s.conn.Close()
}
