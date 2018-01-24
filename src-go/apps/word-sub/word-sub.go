// random(ish) word subscriber (client)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/pubsub"
)

const (
	english = "eng"
	french  = "fra"
	spanish = "esp"
)

const (
	subIDTemplate = "<addr>/word-sub/<pid>"
	subIDFormat   = "%s/word-sub/%d"
)

var subID = flag.String("id", subIDTemplate, "Identifier to be used by the subscriber (no colons!)")
var server = flag.String("server", "localhost", "Name of machine running word publisher")
var port = flag.Int("port", 5678, "Port number to connect to")
var topic = flag.String("lang", "eng", "Language to subscribe to (English='eng', French='fra', Spanish='esp')")
var wordCount = flag.Int("wordcount", 50, "The number of words to get")

func init() {
	log.SetPrefix("word-sub: ")
}

func main() {
	flag.Parse()

	id := getID()
	log.Printf("starting word subscriber app (id: %s)...\n", id)

	d := core.NewDialer("tcp", fmt.Sprintf("%s:%d", *server, *port))
	defer d.Close()

	sub := pubsub.NewSubscriber(id, d)
	defer sub.Close()

	log.Printf("subscribing to %d '%s' words...\n", *wordCount, *topic)
	finished := make(chan bool)
	n := 0

	check.Must(sub.Subscribe(func(m *pubsub.Message) (err error) {
		log.Printf("%s: %s\n", m.Topic, string(m.Body))

		if n++; n == *wordCount {
			log.Println("finishing...")
			finished <- true
		}

		return
	}, *topic))

	// we have nothing else to do but wait...
	<-finished

	log.Println("done.")
}

func getID() string {
	if *subID != subIDTemplate {
		return *subID
	}

	return fmt.Sprintf(subIDFormat, getHostName(), os.Getpid())
}

func getHostName() string {
	host, err := os.Hostname()
	if err != nil {
		// TODO: we could get ip address
		host = "unknown-host"
	}
	return host
}
