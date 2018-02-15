package inputs

var (
	// Hello sent from client to server.
	Hello = &fsm.String{
		Hint:  "Initial message sent from client to server",
		Match: "HELLO",
	}

	// NeedList is sent from client to server.
	NeedList = &fsm.String{
		Hint:  "Initial message sent from server to a primary server",
		Match: "NEEDLIST",
	}

	// Greeting is the first response message sent to either a client or non-primary server.
	Greeting = &GreetingInput{}
	// Any input accepted.
	Any = &fsm.Any{}
)

// GreetingInput describes a valid state machine input of type msgs.Greeting
type GreetingInput struct {
}
