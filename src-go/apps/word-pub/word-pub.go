// random(ish) word publisher (server)
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/pubsub"
)

const (
	topicEnglish = "eng"
	topicFrench  = "fra"
	topicSpanish = "esp"
)

var port = flag.Int("port", 5678, "Port number to listen at")

func init() {
	rand.Seed(time.Now().Unix())
	log.SetPrefix("word-pub: ")
}

func main() {
	log.Printf("starting publisher app (port #%d)...\n", *port)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	listener, err := net.ListenTCP("tcp", addr)
	check.Error(err)
	defer listener.Close()

	pub := pubsub.NewPublisher(listener)
	defer pub.Close()

	for {
		pub.Publish(nextWord())
		time.Sleep(time.Millisecond * 100)
	}
}

func nextWord() (topic string, word []byte) {
	words := dictionary[rand.Intn(len(dictionary))]

	switch rand.Intn(3) {
	case 0:
		topic = topicEnglish
		word = []byte(words.english)
	case 1:
		topic = topicFrench
		word = []byte(words.french)
	case 2:
		topic = topicSpanish
		word = []byte(words.spanish)
	}

	log.Printf("%s: %s\n", topic, string(word))

	return
}
