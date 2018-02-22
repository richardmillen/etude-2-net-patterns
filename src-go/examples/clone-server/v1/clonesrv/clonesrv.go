package main

import (
	"flag"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var primaryAddr = flag.String("primary-server", "", "specifies the address of another active clone server instance.")

var (
	laddr   = flag.String("address", "", "local ip address; an empty value will cause the server to listen on all available unicast and anycast ip addresses.")
	minPort = flag.Int("min-port", 5000, "minimum port number to try to bind to.")
	maxPort = flag.Int("max-port", 5005, "maximum port number to try to bind to.")
)

// clone server states:
var (
	recvClient = &fsm.State{
		Name: "receive requests from clients",
		Events: []*fsm.Event{
			{
				Name:  "hello from client",
				Input: inputs.Hello,
				MoveTo: []*fsm.State{
					greetClient,
				},
			},
		},
	}
	greetClient = &fsm.State{
		Name: "greet new client",
		Events: []*fsm.Event{
			{
				Name:  "greeting to client",
				Input: inputs.Greeting,
				MoveTo: []*fsm.State{
					echoState,
				},
			},
		},
	}
	echoState = &fsm.State{
		Name: "echo server",
		Events: []*fsm.Event{
			{
				Name:  "client echo",
				Input: inputs.Any,
			},
		},
	}
)

// primary server states:
var (
	primary = &fsm.State{
		Name: "primary server",
		Substates: []*fsm.State{
			recvPeer,
			recvClient,
		},
	}
	recvPeer = &fsm.State{
		Name: "receive initial request from peer server",
		Events: []*fsm.Event{
			{
				Name:   "'hi' message from secondary server",
				Input:  inputs.HiBuddy,
				MoveTo: []*fsm.State{greetPeer},
			},
		},
	}
	greetPeer = &fsm.State{
		Name: "greet server",
		Events: []*fsm.Event{
			{
				Name:   "response to peer request",
				Input:  inputs.Greeting,
				MoveTo: []*fsm.State{recvAgain},
			},
		},
	}
	recvAgain = &fsm.State{
		Name: "receive subsequent request from peer server",
		Events: []*fsm.Event{
			{
				Name:   "hi again from secondary server",
				Input:  inputs.HiAgain,
				MoveTo: []*fsm.State{greetPeer},
			},
		},
	}
)

// secondary server states:
var (
	initSecondary = &fsm.State{
		Name: "initial secondary server state",
		Events: []*fsm.Event{
			{
				Name:   "'hi' message sent to primary",
				Input:  inputs.HiBuddy,
				MoveTo: []*fsm.State{recvGreeting},
			},
		},
	}
	recvGreeting = &fsm.State{
		Name: "receive greeting from primary server",
		Events: []*fsm.Event{
			{
				Name:   "greeting message",
				Input:  inputs.Greeting,
				MoveTo: []*fsm.State{secondary},
			},
		},
	}
	secondary = &fsm.State{
		Name: "active secondary server",
		Substates: []*fsm.State{
			recvClient,
			refreshActiveList,
		},
	}
	refreshActiveList = &fsm.State{
		Name: "refresh active server list",
		Substates: []*fsm.State{
			hiAgain,
		},
		Events: []*fsm.Event{
			{
				Name:   "lost connection to primary",
				Input:  &fsm.Error{netx.ErrConnLost},
				MoveTo: []*fsm.State{noPrimary},
			},
		},
	}
	hiAgain = &fsm.State{
		Name: "say hi again",
		Events: []*fsm.Event{
			{
				Name:   "'hi again' from secondary to primary",
				Input:  inputs.HiAgain,
				MoveTo: []*fsm.State{recvUpdate},
			},
		},
	}
	recvUpdate = &fsm.State{
		Name: "receive update from primary server",
		Events: []*fsm.Event{
			{
				Name:   "response to 'hi again' request",
				Input:  inputs.Greeting,
				MoveTo: []*fsm.State{hiAgain},
			},
		},
	}
	noPrimary = &fsm.State{
		Name: "primary server gone",
		Events: []*fsm.Event{
			{
				Name:   "",
				Input:  nil,
				MoveTo: []*fsm.State{},
			},
		},
	}
)

func main() {
	flag.Parse()

	listener, err := netx.ListenTCPRange("tcp", laddr, minPort, maxPort)
	check.Error(err)
	defer listener.Close()

	initState := active
	if *primaryAddr == "" {
		initState = primary
	}

	svc := netx.Service{
		Connector:    listener,
		InitialState: initState,
	}
	defer svc.Close()

	svc.Start()

	for {
		select {
		case r := <-svc.Received():
			r.Write(r.Input())
		case <-svc.Closed():
			log.Println("server closed.")
			return
		}
	}

}
