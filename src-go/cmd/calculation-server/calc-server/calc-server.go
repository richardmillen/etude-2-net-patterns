package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var port = flag.Int("port", 5432, "port number to listen at")

func main() {
	flag.Parse()
	
	log.Println("starting calc-server...")

	l, err := netx.ListenTCP("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	svc := service.New(l)

	recvState := state.New("receiving")
	recvState.Accept(&)

	invalidState := state.New("invalid message")
	invalidState.AddSubstate(recvState)

	svc.Initial(recvState)

	svc.OnRecv(func(f []netx.Frame) error {

	})

	log.Println("calc-server started.")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}
