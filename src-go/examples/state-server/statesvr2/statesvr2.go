package main

import (
	"flag"
	"fmt"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/state-server/input"
)

var port = flag.Int("port", 5432, "port number to listen at.")

func main() {
	flag.Parse()

	firstState := fsm.State{
		Name: "first",
		Events: []fsm.Event{
			{
				Input: input.First,
				Next:  secondState,
			},
		},
	}
	secondState := fsm.State{
		Name: "second",
		Events: []fsm.Event{
			{
				Input: input.Second,
				Next:  thirdState,
			},
		},
	}
	thirdState := fsm.State{
		Name: "third",
		Events: []fsm.Event{
			{
				Input: input.Third,
				Next:  doneState,
			},
		},
	}
	doneState := fsm.State{
		Name: "done",
	}

	listener, err := netx.ListenTCP("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	svc := netx.Service{
		Connector:    listener,
		InitialState: firstState,
		FinalState:   doneState,
	}

	for {
		select {
		case c := <-svc.Connected():
			c.Write([]byte(c.State().Name()))
		case e := <-svc.EnteredState():
			e.Write([]byte(e.State().Name()))
		case <-svc.Closed():
			fmt.Println("service closed.")
			return
		}
	}
}
