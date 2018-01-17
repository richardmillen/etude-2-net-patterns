package pubsub

// data is an internal structure used to contain a message
// from the Publisher as it's sent to each subscription.
type data struct {
	topic   string
	content []byte
}
