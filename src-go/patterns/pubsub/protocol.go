package pubsub

// protocol is the interface to the pub-sub wire protocol.
type protocol interface {
	greet(sub *subscription) error
	sendTo(sub *subscription, topic string, data []byte) error
}
