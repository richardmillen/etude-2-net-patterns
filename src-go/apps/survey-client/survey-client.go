package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

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

	endpoints, err := runSurvey()
	check.Must(err)

	if *service != "echo" {
		fmt.Println("as you're not testing the echo service there's nothing more to do.")
		return
	}

	for _, ep := range endpoints {
		echo(ep)
	}
}

func runSurvey() ([]*disco.Endpoint, error) {
	defer timeThis("survey")()

	surveyor := disco.NewSurveyor(*addr)
	return surveyor.Survey(*service)
}

func echo(endpoint *disco.Endpoint) {
	defer timeThis("echo")()

	conn, err := net.Dial("udp", endpoint.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("sending hello...")
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("receiving response...")
	buf := make([]byte, 8)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("received: '%s'\n", string(buf))
}

func timeThis(message string) func() {
	fmt.Println("started: ", message)
	started := time.Now()
	return func() {
		fmt.Printf("elapsed: %v.\n", time.Since(started))
	}
}
