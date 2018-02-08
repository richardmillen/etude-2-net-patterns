package input

var (
	First = &fsm.String{
		Hint:  "1st message sent from server to client",
		Match: "first",
	}
	Second = &fsm.String{
		Hint:  "2nd message sent from server to client",
		Match: "second",
	}
	Third = &fsm.String{
		Hint:  "3rd message sent from server to client",
		Match: "third",
	}
)
