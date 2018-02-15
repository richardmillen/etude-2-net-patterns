// the clone server is a simple echo server, but with some important differences.
//
// this app can be started with the address of another server instance from which
// it will obtain a list of active instances before it starts listening for client
// requests. then whenever a client connects it will be given a list of known good
// server addresses.
//
// the concept of a 'primary' server instance determines how a list of known good
// server addresses is maintained by any given server instance. 'primary' instances
// i.e. those another clone server from
// which to obtain a list of known good server addresses. this is often (but not
// always) the first one started i.e. the one started without the address of another
// active server.

package main

import (
	"flag"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/clone-server/states"
)

var primaryAddr = flag.String("primary-server", "", "specifies the address of another active clone server instance.")

var (
	laddr   = flag.String("address", "", "local ip address; an empty value will cause the server to listen on all available unicast and anycast ip addresses.")
	minPort = flag.Int("min-port", 5000, "minimum port number to try to bind to.")
	maxPort = flag.Int("max-port", 5005, "maximum port number to try to bind to.")
)

func main() {
	flag.Parse()

	listener, err := netx.ListenTCPRange("tcp", laddr, minPort, maxPort)
	check.Error(err)
	defer listener.Close()

	startState := states.Active
	if *primaryAddr == "" {
		startState = states.Primary
	}

	svc := netx.Service{
		Connector:    listener,
		InitialState: startState,
	}
	defer svc.Close()

	svc.Start()

	for {
		select {
		case r := <-svc.Received():
			r.Write(r.Input())
		case <-svc.Closed():
			log.Println("server closed.")
			return
		}
	}

}
