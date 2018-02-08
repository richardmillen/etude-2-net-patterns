package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/hello-server/input"
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

	sendState.Input(input.Hello, input.Hi)
	sendState.Input(quit).Next(exitState)

	recvState.Input(input.World, input.Error)
	recvState.Substate(sendState)

	dialer := netx.NewDialer("tcp", fmt.Sprintf("%s:%d", *server, *port))
	defer dialer.Close()

	svc := netx.Service{
		Connector:    dialer,
		InitialState: sendState,
		FinalState:   exitState,
	}
	defer svc.Close()

	go func() {
		for {
			select {
			case r := <-svc.Received(input.World):
				fmt.Println("received:", input.World.From(r))
			case r := <-svc.Received(input.Error):
				fmt.Println("server error:", input.Error.From(r))
			case <-svc.Closed():
				fmt.Println("service closed.")
				return
			}
		}
	}()

	userInput := bufio.NewScanner(os.Stdin)

	for svc.IsRunning() {
		fmt.Print("hellocli: ")
		userInput.Scan()
		svc.Write(userInput.Bytes())
	}

}
