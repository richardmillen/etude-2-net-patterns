// random(ish) word subscriber (client)
package main

import (
	"flag"
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
var filter = flag.String("lang", "eng", "Language to subscribe to (English='eng', French='fra', Spanish='esp')")
var countDown = flag.Int("countdown", 100, "The number of words to wait for")

func init() {
	log.SetPrefix("word-sub: ")
}

func main() {
	flag.Parse()

	conn, err := connect(*server)
	utils.CheckError(err)

	sub := pattern.NewSubscriber(conn)

	log.Printf("Subscribing to '%s' words...\n", *filter)
	finished := make(chan bool)

	err = sub.Subscribe(func(b []byte) {
		log.Printf("%s ", string(b))

		if *countDown--; *countDown == 1 {
			finished <- true
		}
	}, *filter)

	// we have nothing else to do but wait...
	<-finished
}

func connect(server string) (net.Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", nil, addr)
}
