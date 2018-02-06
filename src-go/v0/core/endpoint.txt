package core

import (
	"net"
	"os"
	"strings"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

// NewEndpoint constructs a new Endpoint with a unique identifier.
func NewEndpoint(address string) *Endpoint {
	return &Endpoint{UUID: uuid.New(), Addr: address}
}

// NewHostEndpoint constructs a new Endpoint with a unique identifer and local host name.
//
// TODO: think of a better name for this.
func NewHostEndpoint() *Endpoint {
	return &Endpoint{UUID: uuid.New(), Addr: GetEndpointAddress(nil)}
}

// Endpoint represents a service endpoint.
//
// TODO: i'm not convinced this is needed, esp. given that each
// Queue has it's own uuid.
type Endpoint struct {
	UUID uuid.Bytes
	Addr string
}

// GetEndpointAddress returns the hostname:port of the Endpoint.
//
// TODO: think of a better name and location for this helper function.
// HACK: this function is horrible! it's trying to do too much and not very well!
func GetEndpointAddress(addr net.Addr) string {
	host, err := os.Hostname()
	check.Error(err)

	if addr == nil {
		return host
	}

	return strings.Replace(addr.String(), "[::]", host, 1)
}
