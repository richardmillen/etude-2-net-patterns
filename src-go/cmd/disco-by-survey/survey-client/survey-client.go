package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
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

	// ...
}
