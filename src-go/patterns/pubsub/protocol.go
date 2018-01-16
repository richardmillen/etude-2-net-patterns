package pubsub

type protocol interface {
	greet(r *subscription) error
	send(topic string, data []byte) error
}
