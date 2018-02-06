package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/core"
)

var port = flag.Int("port", 5432, "port number to listen at")

func main() {
	log.Println("starting calc-server...")

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	listener, err := core.ListenTCP("tcp", addr)
	check.Error(err)
	defer listener.Close()

	s := core.NewService(listener, &calcServer{})
	defer s.Close()

	s.Connect(func(q *core.Queue) error {
		return core.ErrNoImpl
	})
	s.Error(func(error) error {
		return core.ErrNoImpl
	})
	s.Recv(func(v interface{}) error {
		return core.ErrNoImpl
	})

	s.Start()

	log.Println("calc-server started.")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}
