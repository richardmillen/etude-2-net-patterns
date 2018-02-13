// this application subscribes to a never ending series of live messages
// (random words) sent by a publisher (either 'word-pub' or 'word-lvc').
// the app can subscribe to words in one or more languages (topics);
// 'eng', 'fra' & 'esp'.
//
// if this app subscribes directly to a 'word-pub' instance then it will receive
// messages as & when they are sent by the publisher. on the other hand, if it
// subscribes to a 'word-lvc' instance then it will receive an update (the most
// recent relevant message/s) immediately after a connection is established instead
// of having to wait for the next 'live' message.

package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	english = "eng"
	french  = "fra"
	spanish = "esp"
)

var server = flag.String("server", "localhost", "Name of machine running word publisher")
var port = flag.Int("port", 5678, "Port number to connect to")
var topic = flag.String("lang", "eng", "Language to subscribe to (English='eng', French='fra', Spanish='esp')")
var wordCount = flag.Int("wordcount", 50, "The number of words to get")

func init() {
	log.SetPrefix("word-sub: ")
}

func main() {
	flag.Parse()

	log.Println("starting word subscriber app...")

	d := core.NewDialer("tcp", fmt.Sprintf("%s:%d", *server, *port))
	defer d.Close()

	log.Printf("dialer created (id: %s, addr: %s)...\n", d.EP.UUID, d.EP.Addr)

	sub := netx.NewSubscriber(d, *topic)
	defer sub.Close()

	log.Printf("subscribing to %d '%s' words...\n", *wordCount, *topic)
	finished := make(chan bool)
	n := 0

	sub.Error(func(err error) error {
		log.Println("subscriber error:", err)
		finished <- true
		return err
	})

	sub.Subscribe(func(m *netx.Message) (err error) {
		log.Printf("%s: %s\n", m.Topic, string(m.Body))

		if n++; n == *wordCount {
			log.Println("finishing...")
			finished <- true
		}

		return
	})

	sub.Start()

	// we have nothing else to do but wait...
	<-finished

	log.Println("done.")
}
