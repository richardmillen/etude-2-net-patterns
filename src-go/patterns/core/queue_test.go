package core_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"

	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
)

const testPropName = "abc"

func TestNewQueueID(t *testing.T) {
	q := core.NewQueue(nil, 0)
	defer q.Close()

	if q.ID() == nil {
		t.Error("expected new queue id to be non-nil")
	}
}

func TestNewQueueConn(t *testing.T) {
	expected := &fakeReadWriter{}

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

func TestQueueSendWithNoAvailableCapacity(t *testing.T) {
	q := core.NewQueue(nil, 0)
	defer q.Close()

	err := q.Send("should error")
	if err == nil {
		t.Error("expected error, send called on queue with zero capacity")
	}
}

var noSendRecvTestCases = []struct {
	methodName      string
	methodUnderTest func(*core.Queue) error
	err             error
}{
	{
		methodName: "Send",
		methodUnderTest: func(q *core.Queue) error {
			return q.Send("abc")
		},
	},
	{
		methodName: "Recv",
		methodUnderTest: func(q *core.Queue) error {
			_, err := q.Recv()
			return err
		},
		err: &check.FailedError{},
	},
}

func TestQueueWithNoSendReceiver(t *testing.T) {
	for _, tc := range noSendRecvTestCases {
		t.Run(tc.methodName, func(*testing.T) {
			q := core.NewQueue(nil, 1)
			defer q.Close()

			err := tc.methodUnderTest(q)
			if tc.err != nil {
				switch err.(type) {
				case *check.FailedError:
					// TODO: make this error type more appropriate to the operation.
				case nil:
					t.Errorf("expected error %v from %s, got nil", tc.err, tc.methodName)
				default:
					t.Errorf("expected %T error from %s, got %v", tc.err, tc.methodName, err)
				}
			}

			if err != nil {
				return
			}
			q.Wait()

			err = q.Err()
			if err == nil {
				t.Errorf("expected error; called %s without configuring SendReceiver first.", tc.methodName)
			}
		})
	}
}

// TODO: add test cases.
var sendTestCases = []struct {
	name string
	sr   *fakeSendReceiver
	cap  int
}{
	{
		name: "SingleMessageOK",
		sr: &fakeSendReceiver{
			msgResMap: map[string]*sendResult{
				"hello": &sendResult{
					expectErr: nil,
				},
			},
		},
		cap: 1,
	},
}

func TestQueueSend(t *testing.T) {
	for _, tc := range sendTestCases {
		t.Run(tc.name, func(*testing.T) {
			q := core.NewQueue(nil, tc.cap)
			defer q.Close()

			q.SetSendReceiver(tc.sr)

			for msg, res := range tc.sr.msgResMap {
				if err := q.Send(msg); err != res.expectErr {
					t.Errorf("expected send error %v, got %v", res.expectErr, err)
				}
			}

			q.Wait()

			for msg, res := range tc.sr.msgResMap {
				if !res.sent {
					t.Errorf("message '%s' wasn't sent", msg)
				}
			}
		})
	}
}

// TODO: add table-driven TestQueueRecv().

type fakeReadWriter struct {
	io.ReadWriter
}

type sendResult struct {
	sent      bool
	expectErr error
}

type fakeSendReceiver struct {
	msgResMap map[string]*sendResult
}

func (sr *fakeSendReceiver) Send(q *core.Queue, v interface{}) (err error) {
	s := v.(string)
	if _, ok := sr.msgResMap[s]; !ok {
		return fmt.Errorf("unexpected message '%s' sent to %T", s, sr)
	}

	sr.msgResMap[s].sent = true
	return
}

func (sr *fakeSendReceiver) Recv(q *core.Queue) (interface{}, error) {
	return nil, errors.New("fakeSendReceiver.Recv not implemented")
}
