// TODO: add comments / notes.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

const (
	cmdHelp string = "help"
	cmdQuit string = "quit"
)

var server = flag.String("server", "localhost", "the name/address of the calc-server host")
var port = flag.Int("port", 5432, "the port to connect to on the host")
var calcCount = flag.Int("calc-count", 100, "number of random calculations to perform.")

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()

	dialer, err := netx.NewDialer(fmt.Sprintf("%s:%d", *server, *port))
	check.Error(err)
	defer dialer.Close()

	svc := netx.Service{
		Connector: dialer,
	}

	go func() {
		for {
			select {
			case r := <-svc.Received():
				// TODO: process received.
			case <-svc.Closed():
				fmt.Println("service closed.")
				return
			}
		}
	}()

	for n := 0; n < *calcCount; n++ {
		v1 := rand.Int() + 1
		num.Copy(v1, svc)

		// TODO: send operator, then second operand to service.
	}
}

// TODO: make random.
func getRandomOperator() string {
	return "+"
}
