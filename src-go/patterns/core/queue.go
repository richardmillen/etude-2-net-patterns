package core

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

// DefQueueSize is intended to be used as the default size for connection Queues.
const DefQueueSize = 100

const (
	// PropAddressKey is the key/name of the 'addr' Queue property.
	PropAddressKey = "addr"
)

// QueueGreeter is the interface that wraps the protocol Greet method.
type QueueGreeter interface {
	Greet(q *Queue) error
}

// QueueSender is the interface that wraps the protocol Send method.
type QueueSender interface {
	Send(q *Queue, v interface{}) error
}

// QueueReceiver is the interface that wraps the protocol Recv method.
type QueueReceiver interface {
	Recv(q *Queue) (interface{}, error)
}

// SendReceiver is the interface that groups the Greet, Send and Recv methods.
type SendReceiver interface {
	QueueSender
	QueueReceiver
}

// GreetSendReceiver is the interface that groups the Greet, Send and Recv methods.
type GreetSendReceiver interface {
	QueueGreeter
	SendReceiver
}

// NewQueue constructs a new connection queue.
//
// TODO: pass in rwc, or (more likely) rw?
func NewQueue(rw io.ReadWriter, capacity int) *Queue {
	q := &Queue{
		id: uuid.New(),
		rw: rw,
	}
	q.err = make(chan error, capacity)
	q.props = make(map[string]interface{})
	q.ch = make(chan interface{}, capacity)
	q.quit = make(chan bool, 1)
	q.finished = make(chan bool)
	go q.run()
	return q
}

// Queue handles a subscriber connection on the Publisher.
//
// TODO: finished vs. wgFinished?
// TODO: Send() & Recv() both need to return errors, but so too does the actual
// send operation(?) which happens in a goroutine. this is ugly.
type Queue struct {
	err        chan error
	id         uuid.Bytes
	rw         io.ReadWriter
	sr         SendReceiver
	props      map[string]interface{}
	ch         chan interface{}
	quit       chan bool
	finished   chan bool
	wgSend     sync.WaitGroup
	wgFinished sync.WaitGroup
	//busy       bool
	//busyMutex  sync.Mutex
}

// run is the engine of a connection Queue.
func (q *Queue) run() {
	defer func() {
		log.Println("Queue done.")
		//q.wgFinished.Done()
		q.finished <- true
	}()

	for {
		select {
		case v := <-q.ch:
			q.trySend(v)
		case <-q.quit:
			//log.Println("Queue.run: quit received.")
			return
		}
	}
}

// trySend attempts to send data using the Queues GreetSendReceiver.
func (q *Queue) trySend(v interface{}) {
	defer func() {
		q.wgSend.Done()
	}()

	if q.setError(check.NotNil(q.sr, "queue send-receiver")) {
		return
	}

	q.setError(q.sr.Send(q, v))
}

func (q *Queue) setError(err error) bool {
	if err == nil {
		return false
	}

	select {
	case q.err <- err:
		return true
	default:
		return false
	}
}

// Send is called to pass data to a connection Queue.
func (q *Queue) Send(v interface{}) (err error) {
	q.wgSend.Add(1)

	select {
	case q.ch <- v:
	default:
		err = fmt.Errorf("connection queue '%s' is full", q.id)
		q.wgSend.Done()
	}
	return
}

// Recv is called to receive incoming data.
func (q *Queue) Recv() (v interface{}, err error) {
	//log.Println("Queue.Recv: enter.")

	err = check.NotNil(q.sr, "queue send-receiver")
	if q.setError(err) {
		return
	}

	v, err = q.sr.Recv(q)
	err = GetError(err)
	q.setError(err)

	return
}

// Wait will block until any outstanding send operations complete.
func (q *Queue) Wait() {
	q.wgSend.Wait()
}

// Close waits for the Queue to finish processing messages then kills it.
func (q *Queue) Close() error {
	q.Wait()
	q.quit <- true
	<-q.finished
	return nil
}

// ID is the unique identifier of the queue
//
// TODO: is this even needed?
func (q *Queue) ID() uuid.Bytes {
	return q.id
}

// Conn returns a ReadWriter that represents the connection to the remote endpoint.
//
// TODO: review where this is used and why those areas aren't invoking Send/Recv
// on the Queue directly.
func (q *Queue) Conn() io.ReadWriter {
	return q.rw
}

// Cap returns the capacity of the Queue.
func (q *Queue) Cap() int {
	return cap(q.ch)
}

// Err returns an error from the Queue if one (or more) has occurred.
//
// TODO: exporting Err() whereas function callbacks have been used everywhere else with similar intent.
func (q *Queue) Err() (err error) {
	select {
	case err = <-q.err:
		return err
	default:
	}
	return
}

// SetSendReceiver is called to provide a Queue with a GreetSendReceiver.
func (q *Queue) SetSendReceiver(sr SendReceiver) {
	q.sr = sr
}

// Prop returns the value of a property on the queue.
func (q *Queue) Prop(key string) interface{} {
	return q.props[key]
}

// SetProp is called to set a property on the queue.
func (q *Queue) SetProp(key string, value interface{}) {
	q.props[key] = value
}
