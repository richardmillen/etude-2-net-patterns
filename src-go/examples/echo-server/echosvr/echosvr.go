// the echo server demonstrates the most basic usage of a 'Service' where the 'Service'
// instance uses an implicit (default) 'State' that accepts all input and echoes it back to the client.
//
// n.b. one area that requires special attention is the select case that receives inbound messages
// from the client. the line below is crude at best, but something similar is required in order for the
// server to reply to that specific client connection.
//
// r.Output.Write(r.Input)
//
// in other words calling Service.Write(...) will not suffice as it would be equivalent to doing a multicast
// to all connected clients.

package main

import (
	"flag"
	"fmt"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var port = flag.Int("port", 5432, "port number to listen on.")

func main() {
	flag.Parse()

	listener, err := netx.ListenTCP("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)
	defer listener.Close()

	svc := netx.NewService(listener)
	defer svc.Close()

	for {
		select {
		case r := <-svc.Received():
			r.Output.Write(r.Input)
		case <-svc.Closed():
			fmt.Println("service closed.")
			return
		}
	}
}
