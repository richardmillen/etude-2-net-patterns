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
	QueueSender
	QueueReceiver
}

// NewQueue constructs a new connection queue.
//
// TODO: pass in rwc, or (more likely) rw?
func NewQueue(rw io.ReadWriter, capacity int) *Queue {
	q := &Queue{
		id: uuid.New(),
		rw: rw,
	}
	q.Err = make(chan error, capacity)
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
// TODO: exporting Err whereas function callbacks have been used everywhere else with similar intent.
type Queue struct {
	Err        chan error
	id         uuid.Bytes
	rw         io.ReadWriter
	gsr        GreetSendReceiver
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
		//q.wgFinished.Done()
		q.finished <- true
	}()

	for {
		select {
		case v := <-q.ch:
			err := q.trySend(v)
			if err != nil {
				q.Err <- err
			}
		case <-q.quit:
			log.Println("Queue.run: quit received.")
			return
		}
	}
}

// trySend attempts to send data using the Queues GreetSendReceiver.
func (q *Queue) trySend(v interface{}) (err error) {
	defer func() {
		q.wgSend.Done()
	}()

	err = check.NotNil(q.gsr, "queue greet-send-receiver")
	if err != nil {
		return
	}

	log.Println("queue, message received:", v)
	err = q.gsr.Send(q, v)
	return
}

// ID is the unique identifier of the queue
//
// TODO: is this even needed?
func (q *Queue) ID() uuid.Bytes {
	return q.id
}

// Conn returns a ReadWriter that represents the connection to the remote endpoint.
func (q *Queue) Conn() io.ReadWriter {
	return q.rw
}

// Cap returns the capacity of the Queue.
func (q *Queue) Cap() int {
	return cap(q.ch)
}

// SetGSR is called to provide a Queue with a GreetSendReceiver.
func (q *Queue) SetGSR(gsr GreetSendReceiver) {
	q.gsr = gsr
}

// Prop returns the value of a property on the queue.
func (q *Queue) Prop(key string) interface{} {
	return q.props[key]
}

// SetProp is called to set a property on the queue.
func (q *Queue) SetProp(key string, value interface{}) {
	q.props[key] = value
}

// Send is called to pass data to a connection Queue.
// 
// TODO: should this method pass errors to q.Err rather than returning one?
func (q *Queue) Send(v interface{}) error {
	log.Println("Queue.Send: enter.")
	q.wgSend.Add(1)

	select {
	case q.ch <- v:
		log.Println("Queue.Send: message queued.")
		return nil
	default:
		q.wgSend.Done()
		return fmt.Errorf("connection queue '%s' is full", q.id)
	}
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
