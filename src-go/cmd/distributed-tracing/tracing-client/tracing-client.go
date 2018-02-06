package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var host = flag.String("server", "localhost", "name of log-server hsot.")
var port = flag.Int("port", 5959, "port number to connect to.")
var count = flag.Int("count", 100, "number of times run test function.")
var pause = flag.Bool("pause", false, "whether to pause briefly between test function calls.")

func init() {
	log.SetPrefix(fmt.Sprintf("log-client [%d]: ", os.Getpid()))
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("starting log/trace client (addr: %s)...\n", addr)

	func() {
		for n := 1; n <= *count; n++ {
			msg := fmt.Sprintf("running test #%d...", n)

			log.Println(msg)

			// call foo() ...
		}

		// ...
	}()
}

func foo() {
	if *pause {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}

	switch rand.Intn(10) {
	case 0:
		// alert
	case 1, 2:
		// critical
	case 3, 4, 5:
		baz()
	default:
		bar()
	}
}

func bar() {
	if *pause {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}

	switch rand.Intn(10) {
	case 0:
		// error
	case 1:
		// warning
	case 2, 3, 4:
		// debug
	default:
		baz()
	}
}

func baz() {
	if *pause {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}

	switch rand.Intn(10) {
	case 0, 1, 2, 3, 4:
		// notice
	case 5, 6:
		// debug
	default:
		// info
	}
}
