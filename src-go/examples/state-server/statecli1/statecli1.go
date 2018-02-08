// state client #1 enters a receiving state (opening a connection to a 'state' server),
// receives three specific messages (it doesn't care which order) then the server closes
// the connection and the client exits.

package main

import (
	"flag"
	"fmt"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/state-server/msgs"
)

var server = flag.String("server", "localhost", "server name/address to connect to.")
var port = flag.Int("port", 5432, "server port to connect to.")

func main() {
	flag.Parse()

	recvState := fsm.State{
		Name: "receiving",
		Events: []fsm.Event{
			{Input: msgs.First},
			{Input: msgs.Second},
			{Input: msgs.Third},
		},
	}

	dialer, err := netx.NewDialer(fmt.Sprintf("%s:%d", *server, *port))
	check.Error(err)

	svc := netx.Service{
		Connector:    dialer,
		InitialState: &recvState,
	}

	// would be necessary as constructor not used:
	svc.Start()

	for {
		select {
		case in := <-svc.Received():
			fmt.Println("received:", in)
		case <-svc.Closed():
			fmt.Println("service closed.")
			return
		}
	}
}
