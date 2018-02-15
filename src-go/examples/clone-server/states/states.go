package states

import (
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/clone-server/inputs"
)

// client states:
var (
	Hello = &fsm.State{
		Name: "hello to server",
		Events: []*fsm.Event{
			{
				Name:  "hello request",
				Input: inputs.Hello,
				Next:  []*fsm.State{RecvGreeting},
			},
		},
	}
	RecvGreeting = &fsm.State{
		Name: "greeting from server",
		Events: []*fsm.Event{
			{
				Name:  "hello response",
				Input: inputs.Greeting,
				Next:  []*fsm.State{Echo},
			},
		},
	}
)

// server states:
var (
	Primary = &fsm.State{
		Name: "primary server",
		Substates: []*fsm.State{
			GreetClient,
			GreetServer,
		},
	}
	Active = &fsm.State{
		Name: "active server (secondary)",
		Substates: []*fsm.State{
			GreetClient,
			GetActiveServers,
		},
	}
	GreetClient = &fsm.State{
		Name: "greet new client",
		Events: []*fsm.Event{
			{
				Name:  "hello from client",
				Input: inputs.Hello,
			},
			{
				Name:  "greeting to client",
				Input: inputs.Greeting,
				Next:  []*fsm.State{Echo},
			},
		},
	}
	GetActiveServers = &fsm.State{
		Name: "get active servers",
		Events: []*fsm.Event{
			{
				Name: "",
				Input: 
			},
			{
				Name:  "peer greeting",
				Input: inputs.Greeting,
				Next:  []*fsm.State{Active},
			},
		},
	}
	PrimaryUnavail = &fsm.State{
		Name: "",
		Events: []*fsm.Event{}
	}
	GreetServer = &fsm.State{
		Name: "greet new server",
	}
)

// shared states:
var (
	Echo = &fsm.State{
		Name: "send echo to client",
		Events: []*fsm.Event{
			{
				Name:  "echo response",
				Input: inputs.Any,
				Next:  Recv,
			},
		},
	}
)
