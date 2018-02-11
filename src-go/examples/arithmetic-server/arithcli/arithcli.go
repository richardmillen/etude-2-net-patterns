// the arithmetic client sends a series of basic arithmetical expressions to a remote
// server which returns the result.
//
// each arithmetical expression is sent piece by piece i.e.
// 		1. operand
//		2. operator
//		3. operand
// the server then returns the result.
//
// this implementation only supports very simple arithmetic operations i.e. n+n, n/n etc.
// where 'n' is a 32-bit float.
//
// n.b. type safety is moot (if not misplaced) at the point where data is passed into the
// Server because validation (and serialisation) would be performed within the Service by
// the current State, or more accurately by the fsm.Input on the associated event (fsm.Event).
// so code such as the following which copies values to the Service using the relevant
// input types resembles what would happen within a simple call to Service.Send()/.Execute()
// (or whatever the API ends up looking like):
//
// msgs.Num.Copy(a, svc)
// msgs.Op.Copy(op, svc)
// msgs.Num.Copy(b, svc)
//

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

	svc := netx.NewService(dialer)
	defer svc.Close()

	go func() {
		for {
			select {
			case r := <-svc.Received():
				fmt.Println("received:", r.Input)
			case <-svc.Closed():
				fmt.Println("service closed.")
				return
			}
		}
	}()

	for n := 0; n < *calcCount; n++ {
		a := rand.Int() + 1
		op := getRandomOperator()
		b := rand.Int() + 1

		svc.Send(a)
		svc.Send(op)
		svc.Send(b)
	}
}

func getRandomOperator() string {
	switch rand.Intn(3) {
	case 0:
		return "*"
	case 1:
		return "/"
	case 2:
		return "+"
	case 3:
		return "-"
	}
}
