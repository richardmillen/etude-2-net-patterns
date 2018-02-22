package inputs

var (
	// HelloServer sent from client to server.
	HelloServer = &fsm.String{
		Hint:  "Initial message sent from client to server",
		Match: "HELLO",
	}

	// HiPrimary is sent from a secondary server to a primary server.
	HiPrimary = &HiPrimaryInput{}

	// HiAgain is sent from a secondary server to a primary server.
	HiAgain = &fsm.String{
		Hint:  "Message requesting active server list",
		Match: "HIAGAIN",
	}

	// Greeting is the response message sent to either a client or secondary server.
	Greeting = &GreetingInput{}

	// Any input accepted.
	Any = &fsm.Any{}
)

// HiPrimaryInput describes a valid state machine input of type msgs.HiPrimary
type HiPrimaryInput struct {
}

// GreetingInput describes a valid state machine input of type msgs.Greeting
type GreetingInput struct {
}
