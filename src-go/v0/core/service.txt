package core

import (
	"log"
	"sync"
)

// NewService constructs a new Service.
func NewService(c Connector, gsr GreetSendReceiver) *Service {
	s := &Service{
		connector: c,
		gsr:       gsr,
	}

	s.connector.OnConnect(s.onNewConn)
	s.ch = make(chan interface{}, s.connector.QueueSize())
	s.quit = make(chan bool, 2)

	return s
}

// Service sends and receives messages from a remote Service endpoint.
type Service struct {
	connector Connector
	gsr       GreetSendReceiver
	connFunc  ConnectFunc
	recvFunc  RecvFunc
	errFunc   ErrorFunc
	ch        chan interface{}
	quit      chan bool
	wgSend    sync.WaitGroup
	wgFinish  sync.WaitGroup
}

//finished  chan bool

// Start is called to activate a Service after it has been configured.
//
// Messages cannot be received from a remote endpoint until a Service
// has been Start'ed.
//
// Note that the select case/default as opposed to select case/case
// where the latter includes the quit channel. this is to ensure the
// Service channel is flushed before responding to the quit channel.
// put another way, several messages could be queued in s.ch then the
// application could close, causing an event on s.quit. this would
// mean that anything in the s.ch queue would be lost.
// n.b. if this behaviour is desirable then it should still be possible
// by configuring the connector to quit before the Service. GetQueues
// could be made to return nil for instance.
//
// Refer to the language spec for furter info on select case/case vs
// select case/default:
// https://golang.org/ref/spec#Select_statements
//
// TODO: handle connection errors by retrying.
// TODO: should we report queue errors to the consumer?
func (s *Service) Start() {
	go func() {
		s.wgFinish.Add(1)
		defer s.wgFinish.Done()

		err := s.connector.Open(s.gsr)
		if s.raiseIfError(err) {
			return
		}

		go func() {
			s.wgFinish.Add(1)
			defer s.wgFinish.Done()

			for {
				log.Println("waiting for inbound message...")
				err := RecvQueues(s.connector, s.onRecv, s.onError)
				if err != nil {
					log.Println(err)
					s.quit <- true
					return
				}
			}
		}()

		for {
			select {
			case m := <-s.ch:
				SendToQueues(s.connector, &m)
			case <-s.quit:
				CloseQueues(s.connector)
				return
			}
		}
	}()
}

func (s *Service) onNewConn(q *Queue) error {
	if s.connFunc == nil {
		return nil
	}
	return s.connFunc(q)
}

// onRecv forwards a received message to the Services RecvFunc if it's configured.
func (s *Service) onRecv(v interface{}) error {
	if s.recvFunc == nil {
		return nil
	}
	return s.recvFunc(v)
}

// onError forwards a Queue receive error to the Services ErrorFunc if it's configured.
func (s *Service) onError(err error) error {
	if s.errFunc == nil {
		return nil
	}
	return s.errFunc(err)
}

func (s *Service) raiseIfError(err error) bool {
	if err == nil {
		return false
	}

	s.onError(err)
	return true
}

// Send is called to send data to connected endpoint(s).
func (s *Service) Send(v interface{}) error {
	SendToQueues(s.connector, v)
	return nil
}

// Write writes len(p) bytes from p to connected endpoint(s).
func (s *Service) Write(p []byte) (n int, err error) {
	SendToQueues(s.connector, p)
	return 0, nil

	// TODO: reinstate the channel-based logic.
	/*s.wgSend.Add(1)

	select {
	case s.ch <- v:
		fmt.Println("Service.Write: message sent to channel.")
		return nil
	default:
		s.wgSend.Done()
		return errors.New("service queue full")
	}*/
}

// Connect is called to configure the ConnectFunc of a Service.
//
// TODO: this is flawed as connection could well have been established
// before this method is called.
func (s *Service) Connect(connFunc ConnectFunc) {
	s.connFunc = connFunc
}

// Error is called to configure the ErrorFunc of a Service,
// which is executed if a runtime error occurs while subscribing.
func (s *Service) Error(errFunc ErrorFunc) {
	s.errFunc = errFunc
}

// Recv is called to configure the RecvFunc of a Service,
// which is executed every time data is received from a remote Service.
//
// TODO: cater for multiple calls and from multiple goroutines.
// TODO: test Close() then Recv()
func (s *Service) Recv(recvFunc RecvFunc) {
	s.recvFunc = recvFunc
}

// Close is called to stop and invalidate the Service.
func (s *Service) Close() error {
	s.wgSend.Wait()
	s.quit <- true
	s.wgFinish.Wait()
	return nil
}
