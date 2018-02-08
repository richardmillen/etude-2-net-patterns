// the 'state' server example #1 sends a sequence of messages to clients immediately
// after they connect, then upon reaching the end of the sequence closes the connection.
// an 'extra' state (and resulting message) can be enabled via a command line flag that
// will test the clients response to unexpected messages.

package main

import (
	"flag"
	"fmt"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/state-server/msgs"
)

var port = flag.Int("port", 5432, "port number to listen on.")
var includeExtra = flag.Bool("send-extra", false, "enable this to send an extra message to the client and see how it responds.")

func main() {
	flag.Parse()

	firstState := fsm.NewState("first")
	secondState := fsm.NewState("second")
	thirdState := fsm.NewState("third")
	extraState := fsm.NewState("extra")
	doneState := fsm.NewState("done")

	firstState.Event(msgs.First).MoveTo(secondState)
	secondState.Event(msgs.Second).MoveTo(thirdState)

	if *includeExtra {
		thirdState.Event(msgs.Third).MoveTo(extraState)
		extraState.Event(&fsm.Any{}).MoveTo(doneState)
	} else {
		thirdState.Event(msgs.Third).MoveTo(doneState)
	}

	go func() {
		for {
			select {
			case e := <-firstState.Entered():
				e.State.Write([]byte(e.State().Name()))
			case e := <-secondState.Entered():
				e.State.Write([]byte(e.State().Name()))
			case e := <-thirdState.Entered():
				e.State.Write([]byte(e.State().Name()))
			case e := <-extraState.Entered():
				e.State.Write([]byte("this is an extra message just to see how the client reacts!"))
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
