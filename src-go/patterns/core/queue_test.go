package core_test

import (
	"io"
	"testing"

	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
)

const testPropName = "abc"

type fakeProto struct {
	core.GreetSendReceiver
}

type fakeConn struct {
	io.Reader
	io.Writer
}

func TestNewQueueID(t *testing.T) {
	q := core.NewQueue(nil, 0)
	defer q.Close()

	if q.ID() == nil {
		t.Error("expected new queue id to be non-nil")
	}
}

func TestNewQueueConn(t *testing.T) {
	expected := &fakeConn{}

	q := core.NewQueue(expected, 0)
	defer q.Close()

	if q.Conn() != expected {
		t.Errorf("expected %v, got %v", expected, q.Conn())
	}
}

func TestNewQueueCap(t *testing.T) {
	expected := 10

	q := core.NewQueue(nil, expected)
	defer q.Close()

	if q.Cap() != expected {
		t.Errorf("expected %d, got %d", expected, q.Cap())
	}
}

func TestQueueEmptyProp(t *testing.T) {
	q := core.NewQueue(nil, 0)
	defer q.Close()

	if q.Prop(testPropName) != nil {
		t.Errorf("expected %v, got %v", nil, q.Prop(testPropName))
	}
}

func TestQueueSetProp(t *testing.T) {
	expected := "value"

	q := core.NewQueue(nil, 0)
	defer q.Close()

	q.SetProp(testPropName, expected)

	if q.Prop(testPropName) != expected {
		t.Errorf("expected %v, got %v", expected, q.Prop(testPropName))
	}
}

func TestQueueSendWithNoGSR(t *testing.T) {
	q := core.NewQueue(nil, 1)
	defer q.Close()

	err := q.Send("hello")
	if err != nil {
		t.Errorf("send failed with error %v", err)
	}
	q.Wait()

	select {
	case <-q.Err:
	default:
		t.Error("expected error; q.SetGSR wasn't called prior to q.Send.")
	}
}

// TODO: make table-driven tests.
func TestQueueSend(t *testing.T) {
	conn := &fakeConn{}
	capacity := 10

	q := core.NewQueue(conn, capacity)
	defer q.Close()

}
