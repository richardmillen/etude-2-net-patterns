package internal

// Msg is an internal structure used to contain a message from a Publisher.
type Msg struct {
	Topic string
	Body  []byte
}
