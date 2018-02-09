// the 'arithmetic' server accepts a sequence of messages from a client then when
// the sequence is complete performs the arithmetical operation and returns the
// result to the client.
//
// this example would support only a single basic arithmetic operation on two 32-bit
// floating point values i.e. 3.0f+7.0f, 9.0f/9.0f and so on.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var port = flag.Int("port", 5432, "port number to listen at")

var (
	errorState = &fsm.State{
		Name: "error",
		Events: []*fsm.Event{
			{
				Input: any,
			},
		},
		Substates: []*fsm.State{
			numState,
			opState,
			calcState,
		},
	}
)

// calculation represents an ongoing calculation being performed for a client.
// it embeds a hypothetical netx.Conn type (needs a better name!) that is the
// connection to the client.
// the purpose of this type is to maintain state until the calculation can be
// performed and the result returned to the client.
type calculation struct {
	netx.Conn
	operands []float32
	operator *operator
}

// newCalculation is required in order for the Service's Listener to construct
// a 'calculation' rather than a basic netx.Conn.
// see the call to netx.Listener.SetConstructor() below.
func newCalculation() *netx.Conn {
	return &calculation{
		operands: make([]float32, 0, 2),
	}
}

func main() {
	flag.Parse()

	listener, err := netx.ListenTCP("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)
	defer listener.Close()

	listener.SetConstructor(newCalculation)

	svc := &netx.Service{
		Connector:    listener,
		InitialState: numState,
		FinalState:   calcState,
	}
	defer svc.Close()

	for {
		select {
		case r := <-svc.ReceivedInput(num):
			calc := r.State.(calculation)
			calc.operands = append(calc.operands, num.From(r))
		case r := <-svc.ReceivedInput(op):
			calc := r.State.(calculation)
			calc.operators = op.From(r)
		case r := <-svc.ReceivedInput(any):
			log.Println("received:", r.Input)
			e.State.Write([]byte(fmt.Sprintf("invalid message: %v", r.Input)))
		case e := <-svc.EnteredState(calcState):
			calc := e.State.(calculation)
			result := calc.operators[0].oper(calc.operands[0], calc.operands[1])
			num.Copy(result, calc)
		case <-svc.Closed():
			log.Println("service closed.")
			return
		}
	}
}
