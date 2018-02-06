package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var port = flag.Int("port", 5432, "port number to listen at")

func main() {
	log.Println("starting calc-server...")

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	listener, err := netx.ListenTCP("tcp", addr)
	check.Error(err)
	defer listener.Close()

	s := netx.NewService(listener, &calcServer{})
	defer s.Close()

	log.Println("calc-server started.")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}
