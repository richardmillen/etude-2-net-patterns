// 'state' client #2 receives three specific messages from a 'state' server.
// the messages must be in the correct order, or an error is displayed and
// the client aborts.

package main

import (
	"flag"
	"fmt"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/state-server/msgs"
)

var server = flag.String("server", "localhost", "server name/address to connect to.")
var port = flag.Int("port", 5432, "server port to connect to.")

func main() {
	flag.Parse()

	states := []fsm.State{
		{
			Name: "first",
			Transitions: []*fsm.Transition{
				{
					Input:      msgs.First,
					TargetName: "second",
				},
			},
		},
		{
			Name: "second",
			Transitions: []*fsm.Transition{
				{
					Input:      msgs.Second,
					TargetName: "third",
				},
			},
		},
		{
			Name: "third",
			Accepts: []fsm.Input{
				msgs.Third,
			},
		},
		{
			Name: "finished",
		},
	}

	invalidState := fsm.State{
		Name:      "invalid",
		Accepts:   []fsm.Input{&fsm.Any{}},
		Substates: states,
	}

	go func() {
		select {
		case msg := <-invalidState.Received():
			invalidState.Machine.Abort(fmt.Errorf("invalid message received: %s", msg))
		case <-invalidState.Exited():
		}
	}()

	dialer, err := netx.NewDialer(fmt.Sprintf("%s:%d", *server, *port))
	check.Error(err)
	defer dialer.Close()

	svc := netx.Service{
		Connector:    dialer,
		InitialState: states[0],
		FinalState:   states[len(states)-1],
	}
	defer svc.Close()
	svc.Start()

	for {
		select {
		case msg := <-svc.Received():
			fmt.Println("received:", msg)
		case err := <-svc.Closed():
			if err != nil {
				fmt.Println("service error:", err)
			}
			return
		}
	}
}
