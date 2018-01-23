package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/apps/services/echo"
	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/disco"
)

var service = flag.String("service", "echo", "name of service to look for.")
var addr = flag.String("addr", "localhost", "name or ip address. can be broadcast (IP v4 only), multicast or unicast address.")

func main() {
	flag.Parse()

	defer func() {
		fmt.Println("done.")
	}()

	finished := make(chan bool)

	surveyor := disco.NewSurveyor(*addr)
	defer surveyor.Close()

	check.Must(surveyor.Survey(func(endpoint *disco.Endpoint) error {
		defer timeThis("echo (survey response)")()

		conn, err := net.Dial("tcp", endpoint.Addr)
		check.Error(err)
		defer conn.Close()

		check.Must(echo.Send(conn, "hello"))

		rep, err := echo.Recv(conn)
		check.Error(err)

		fmt.Printf("received: '%s'\n", string(rep))

		finished <- true
		return disco.ErrEndSurvey
	}, time.Second, *service))

	<-finished
}

func timeThis(message string) func() {
	fmt.Println("started: ", message)
	started := time.Now()
	return func() {
		fmt.Printf("elapsed: %v.\n", time.Since(started))
	}
}
