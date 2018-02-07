package hello

import "flag"

var port = flag.Int("port", 5432, "port number to connect at")

var (
	Hello = &fsm.String{
		Hint:  "'hello' message to server",
		Match: "hello",
	}
	Hi = &fsm.String{
		Hint:  "'hi' message to server",
		Match: "hi",
	}
	World = &fsm.String{
		Hint:  "'world' response from server",
		Match: "world",
	}
	Error = &fsm.String{
		Hint:  "error response from server",
		Match: "invalid request",
	}
)
