package msgs

var (
	// Hello message
	Hello = &fsm.String{
		Hint:  "'hello' message to server",
		Match: "hello",
	}
	// Hi message
	Hi = &fsm.String{
		Hint:  "'hi' message to server",
		Match: "hi",
	}
	// World message
	World = &fsm.String{
		Hint:  "'world' response from server",
		Match: "world",
	}
	// Error message
	Error = &fsm.String{
		Hint:  "error response from server",
		Match: "invalid request",
	}
)
