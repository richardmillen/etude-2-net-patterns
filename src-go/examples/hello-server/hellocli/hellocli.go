// the hello world client allows the user to enter string messages to be sent to
// a hello word server. however only "hi" and "hello" are valid request messages.
// entering "quit" will shut down the client service and terminate the application.
//
// note that the client Service must be started somehow in order to spawn background
// goroutines etc. this example does so explicitly with a call to Service.Start().

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/hello-server/msgs"
)

var server = flag.String("server", "localhost", "server name/address to connect to.")
var port = flag.Int("port", 5432, "server port number to connect to.")

var quit = &fsm.String{
	Hint:  "'quit' command",
	Match: "quit",
}

func main() {
	flag.Parse()

	sendState := fsm.NewState("send")
	recvState := fsm.NewState("receive")
	exitState := fsm.NewState("exit")

	sendState.Input(msgs.Hello, msgs.Hi)
	sendState.Input(quit).Next(exitState)

	recvState.Input(msgs.World, msgs.Error)
	recvState.Substate(sendState)

	dialer := netx.NewDialer("tcp", fmt.Sprintf("%s:%d", *server, *port))
	defer dialer.Close()

	svc := netx.NewService{
		Connector:    dialer,
		InitialState: sendState,
		FinalState:   exitState,
	}
	defer svc.Close()

	go func() {
		for {
			select {
			case r := <-svc.Received(msgs.World):
				fmt.Println("received:", msgs.World.From(r))
			case r := <-svc.Received(msgs.Error):
				fmt.Println("server error:", msgs.Error.From(r))
			case <-svc.Closed():
				fmt.Println("service closed.")
				return
			}
		}
	}()

	svc.Start()

	userInput := bufio.NewScanner(os.Stdin)

	for svc.IsRunning() {
		fmt.Print("hellocli: ")
		userInput.Scan()
		svc.Write(userInput.Bytes())
	}

}
