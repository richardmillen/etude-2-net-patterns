// random(ish) word publisher (server)
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/pattern"
	"github.com/richardmillen/etude-2-net-patterns/src-go/utils"
)

const (
	filterEnglish = "eng"
	filterFrench  = "fra"
	filterSpanish = "esp"
)

var port = flag.Int("port", 5678, "Port number to listen on")

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	utils.CheckError(err)

	listener, err := net.ListenTCP("tcp", addr)
	utils.CheckError(err)

	pub := pattern.NewPublisher(listener)
	defer pub.Close()

	for {
		pub.Publish(next())
	}
}

func next() (filter string, word []byte) {
	words := dictionary[rand.Intn(len(dictionary))]

	switch rand.Intn(3) {
	case 0:
		filter = filterEnglish
		word = []byte(words.english)
	case 1:
		filter = filterFrench
		word = []byte(words.french)
	case 2:
		filter = filterSpanish
		word = []byte(words.spanish)
	}
	return
}
