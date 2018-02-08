package main

import (
	"flag"
	"fmt"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/state-server/input"
)

var port = flag.Int("port", 5432, "port number to listen on.")

func main() {
	flag.Parse()

	firstState := fsm.NewState("first")
	secondState := fsm.NewState("second")
	thirdState := fsm.NewState("third")
	doneState := fsm.NewState("done")

	firstState.Event(input.First).MoveTo(secondState)
	secondState.Event(input.Second).MoveTo(thirdState)
	thirdState.Event(input.Third).MoveTo(doneState)

	go func() {
		for {
			select {
			case e := <-firstState.Entered():
				e.State.Write(e.State().Name())
			case e := <-secondState.Entered():
				e.State.Write(e.State().Name())
			case e := <-thirdState.Entered():
				e.State.Write(e.State().Name())
			}
		}
	}()

	listener, err := netx.ListenTCP("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	svc := netx.NewService(listener)

	svc.SetInitialState(firstState)
	svc.SetFinalState(doneState)

	svc.Run()
}
