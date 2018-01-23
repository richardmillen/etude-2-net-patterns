package main

import (
	"flag"
	"log"
	"math/rand"
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
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()
	delay()

	defer diags.Start("echo")()

	log.Println("starting survey client...")

	surveyor := disco.NewSurveyor(*addr)
	defer surveyor.Close()

	finished := make(chan bool)

	check.Must(surveyor.Survey(func(addr string) error {
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

// delay was added as a quick hack to enable higher volumes.
func delay() {
	log.Println("waiting for random (short) duration...")

	ms := rand.Intn(2000) + 1000
	time.Sleep(time.Millisecond * time.Duration(ms))

	log.Println("resuming...")
}
