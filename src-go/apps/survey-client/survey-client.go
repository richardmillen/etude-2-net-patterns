package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/apps/services/echo"
	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/diags"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/disco"
)

var service = flag.String("service", "echo", "name of service to look for.")
var addr = flag.String("addr", "localhost", "name or ip address. can be broadcast (IP v4 only), multicast or unicast address.")
var echoText = flag.String("echo-text", "hello", "text to be sent to echo server.")

func init() {
	log.SetPrefix("survey-client: ")
}

func main() {
	flag.Parse()

	defer func() {
		log.Println("done.")
	}()

	log.Println("starting survey client...")

	finished := make(chan bool)

	surveyor := disco.NewSurveyor(*addr)
	defer surveyor.Close()

	check.Must(surveyor.Survey(func(addr string) error {
		defer diags.Start("echo")()

		conn, err := net.Dial("tcp", addr)
		check.Error(err)
		defer conn.Close()

		log.Printf("sending '%s'\n", *echoText)

		check.Must(echo.Send(conn, *echoText))

		rep, err := echo.Recv(conn)
		check.Error(err)

		log.Printf("received: '%s'\n", rep)

		finished <- true
		return disco.ErrEndSurvey
	}, time.Second, *service))

	<-finished
}
