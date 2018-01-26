package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/core"
	"github.com/richardmillen/etude-2-net-patterns/src-go/patterns/pubsub"
	"github.com/richardmillen/etude-2-net-patterns/src-go/services/logger"
)

var host = flag.String("server", "localhost", "name of log-server hsot.")
var port = flag.Int("port", 5959, "port number to connect to.")
var count = flag.Int("count", 100, "number of times run test function.")

var (
	alerts    = 0
	criticals = 0
	errors    = 0
	warnings  = 0
	notices   = 0
)

func init() {
	log.SetPrefix(fmt.Sprintf("log-client [%d]: ", os.Getpid()))
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("starting log/trace client (addr: %s)...\n", addr)

	d := core.NewDialer("tcp", addr)
	defer d.Close()

	pub := pubsub.NewPublisher(d)
	defer pub.Close()

	for n := 0; n < *count; n++ {
		foo(pub)
	}
}

func foo(pub *pubsub.Publisher) {
	log := logger.Start(pub, "foo")
	defer log.Close()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

	switch rand.Intn(10) {
	case 0:
		log.Printf(logger.Alert, "alert message %d", alerts)
		alerts++
	case 1, 2:
		log.Printf(logger.Critical, "critical message %d", criticals)
		criticals++
	case 3, 4, 5:
		baz(log)
	default:
		bar(log)
	}
}

func bar(opLog *logger.Logger) {
	log := logger.StartChild(opLog, "bar")
	defer log.Close()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

	switch rand.Intn(10) {
	case 0:
		log.Printf(logger.Error, "error message %d", errors)
		errors++
	case 1:
		log.Printf(logger.Warning, "warning message %d", warnings)
		warnings++
	case 2, 3, 4:
		log.Print(logger.Debug, "debug message")
	default:
		baz(opLog)
	}
}

func baz(opLog *logger.Logger) {
	log := logger.StartChild(opLog, "baz")
	defer log.Close()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

	switch rand.Intn(10) {
	case 0, 1, 2, 3, 4:
		log.Printf(logger.Notice, "notice message %d", notices)
		notices++
	case 5, 6:
		log.Print(logger.Debug, "debug message")
	default:
		log.Print(logger.Info, "info message")
	}
}
