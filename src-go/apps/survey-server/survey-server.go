package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/richardmillen/etude-2-net-patterns/src-go/apps/services/echo"
	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/disco"
)

// change default port to 0 (zero) to use ephemeral port.
var port = flag.Int("echo-port", 5858, "port number to listen for echo requests at.")

func main() {
	flag.Parse()

	defer func() {
		fmt.Println("done.")
	}()

	log.Printf("starting discoverable echo server (port #%d)...\n", *port)

	echo := startEchoServer()
	defer echo.Close()

	candidate := startCandidate(echo)
	defer candidate.Close()

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	<-ctrlC
}

func startEchoServer() *echo.Server {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	listener, err := net.ListenTCP("tcp", addr)
	check.Error(err)
	defer listener.Close()

	return echo.NewServer(listener)
}

// HACK: this assumes that the echo server address will remain the same throughout the session.
func startCandidate(s *echo.Server) *disco.Candidate {
	candidate := disco.NewCandidate()
	candidate.AddService(echo.ServiceName, s.Addr())
	check.Must(candidate.Open())
	return candidate
}
