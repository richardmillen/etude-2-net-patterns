package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/pubsub"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var host = flag.String("server", "localhost", "name of log-server hsot.")
var port = flag.Int("port", 5959, "port number to connect to.")
var events = flag.Int("events", 100, "number of random events to log.")

func init() {
	log.SetPrefix(fmt.Sprintf("log-client [%d]: ", os.Getpid()))
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	check.Error(err)
	defer conn.Close()

	pub := pubsub.NewPublisher()
	defer pub.Close()

	for {
		pub.Publish(nextEvent())
		time.Sleep(time.Millisecond * 100)
	}
}

func nextEvent() (string, []byte]) {

}
