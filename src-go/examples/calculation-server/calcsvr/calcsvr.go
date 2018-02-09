// TODO: add comments & notes.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var port = flag.Int("port", 5432, "port number to listen at")

const opersExpr = `[\+\-/\*]`

type operator struct {
	fsm.Regex
}

var (
	num = &fsm.Float{
		Hint: "numeric operand",
	}
	op     = &operator{}
	equals = &fsm.String{
		Hint:  "equals sign",
		Match: "=",
	}
	any = &fsm.Any{}
)

var (
	numState = &fsm.State{
		Name: "number",
		Events: []*fsm.Event{
			{
				Input: num,
				MoveTo: []*fsm.State{
					opState,
					calcState,
				},
			},
		},
	}
	opState = &fsm.State{
		Name: "operator",
		Events: []*fsm.Event{
			{
				Input:  op,
				MoveTo: []*fsm.State{numState},
			},
		},
	}
	calcState = &fsm.State{
		Name: "calculate",
		Events: []*fsm.Event{
			{
				Input: num,
			},
		},
	}
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

type calculation struct {
	netx.Conn
	operands  []float32
	operators []*operator
}

func newCalculation() *netx.Conn {
	return &calculation{
		operands:  make([]float32, 0, 2),
		operators: make([]*operator, 0, 1),
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
			calc.operators = append(calc.operators, op.From(r))
		case e := <-svc.EnteredState(calcState):
			calc := e.State.(calculation)
			result := calc.operators[0].Oper(calc.operands[0], calc.operands[1])
			num.Copy(result, calc)
		case r := <-svc.ReceivedInput(any):
			log.Println("received:", r.Input)
			e.State.Write([]byte(fmt.Sprintf("invalid message: %v", r.Input)))
		case <-svc.Closed():
			log.Println("service closed.")
			return
		}
	}
}
