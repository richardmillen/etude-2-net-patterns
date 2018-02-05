package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/core"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":5432")
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

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}
