// random(ish) word subscriber (client)
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/pattern"
	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
)

const (
	english = "eng"
	french  = "fra"
	spanish = "esp"
)

var server = flag.String("server", "localhost", "Name of machine running word publisher")
var port = flag.Int("port", 5678, "Port number to connect to")
var filter = flag.String("lang", "eng", "Language to subscribe to (English='eng', French='fra', Spanish='esp')")
var wordCount = flag.Int("wordcount", 100, "The number of words to get")

func init() {
	log.SetPrefix("word-sub: ")
}

func main() {
	flag.Parse()

	conn, err := connect(*server, *port)
	utils.CheckError(err)

	sub := pattern.NewSubscriber(conn)

	log.Printf("Subscribing to '%s' words...\n", *filter)
	finished := make(chan bool)
	n := 0

	err = sub.Subscribe(func(b []byte) {
		log.Printf("%s ", string(b))

		if n++; n == *wordCount {
			finished <- true
		}
	}, *filter)

	// we have nothing else to do but wait...
	<-finished
}

func connect(server string, port int) (net.Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", nil, addr)
}
