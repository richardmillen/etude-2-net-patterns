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
// TODO: finish this example.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/examples/arithmetic-server/msgs"
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

		msgs.Num.Copy(a, svc)
		msgs.Op.Copy(op, svc)
		msgs.Num.Copy(b, svc)
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
