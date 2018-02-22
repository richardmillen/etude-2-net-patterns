//

package main

import (
	"flag"

	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/clone-server/inputs"
)

var (
	helloState = &fsm.State{
		Name: "hello to server",
		Events: []*fsm.Event{
			{
				Name:   "hello request",
				Input:  inputs.Hello,
				MoveTo: []*fsm.State{greetState},
			},
		},
	}
	greetState = &fsm.State{
		Name: "greeting from server",
		Events: []*fsm.Event{
			{
				Name:   "hello response",
				Input:  inputs.Greeting,
				MoveTo: []*fsm.State{echoState},
			},
		},
	}
	echoState = &fsm.State{
		Name: "echo",
		Events: []*fsm.Event{
			{
				Name:  "echo message",
				Input: inputs.Any,
			},
			{
				Name:   "lost connection to server",
				Input:  &fsm.Error{netx.ErrConnLost},
				MoveTo: []*fsm.State{disconState},
			},
		},
	}
	disconState = &fsm.State{
		Name: "disconnected",
		Events: []*fsm.Event{
			{
				Name:   "",
				Input:  inputs.Any,
				MoveTo: []*fsm.State{helloState},
			},
		},
	}
)

func main() {
	flag.Parse()

}
