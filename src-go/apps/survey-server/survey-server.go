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

// default port is 0 (zero) in order to use ephemeral port.
var port = flag.Int("echo-port", 0, "port number to listen for echo requests at.")

func init() {
	log.SetPrefix("survey-server: ")
}

func main() {
	flag.Parse()

	defer func() {
		log.Println("done.")
	}()

	log.Println("starting survey server...")

	echo := startEchoServer()
	defer echo.Close()

	candidate := startCandidate(echo)
	defer candidate.Close()

	log.Println("survey server started.")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}

func startEchoServer() *echo.Server {
	log.Printf("starting echo server on %s...\n", echoPortDesc())

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	listener, err := net.ListenTCP("tcp", addr)
	check.Error(err)

	return echo.NewServer(listener)
}

// HACK: this assumes that the echo server address will remain the same throughout the session.
func startCandidate(s *echo.Server) *disco.Candidate {
	log.Println("starting service candidate...")

	candidate := disco.NewCandidate()
	candidate.AddService(echo.ServiceName, s.Addr())
	check.Must(candidate.Open())
	return candidate
}

func echoPortDesc() string {
	if *port == 0 {
		return "ephemeral port"
	}
	return fmt.Sprintf("port %d", *port)
}
