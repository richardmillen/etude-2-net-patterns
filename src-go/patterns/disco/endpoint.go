package disco

import "github.com/richardmillen/etude-2-net-patterns/src-go/uuid"

// NewEndpoint constructs a new Endpoint with unique identifier.
func NewEndpoint(addr string) *Endpoint {
	return &Endpoint{UUID: uuid.New(), Addr: addr}
}

// Endpoint represents a discovery endpoint.
type Endpoint struct {
	UUID uuid.Bytes
	Addr string
}
