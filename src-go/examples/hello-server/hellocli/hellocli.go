package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/hello-server/hello"
)

var server = flag.String("server", "localhost", "server name/address to connect to.")

var quit = &fsm.String{
	Hint:  "'quit' command",
	Match: "quit",
}

func main() {
	flag.Parse()

	sendState := fsm.NewState("send")
	recvState := fsm.NewState("receive")
	exitState := fsm.NewState("exit")

	sendState.Accept(hello.Hello, hello.Hi, quit)
	recvState.Accept(hello.World, hello.Error)

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
			case r := <-svc.Received(hello.World):
				fmt.Println("received:", hello.World.From(r))
			case r := <-svc.Received(hello.Error):
				fmt.Println("server error:", hello.Error.From(r))
			case r := <-svc.Received(quit):
				// TODO: transition to exitState => how to implement transitions?
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
