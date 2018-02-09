// the echo client demonstrates the most basic usage of a 'Service' where the
// 'Service' instance uses an implicit (default) 'State' that accepts all input
// entered by the user or sent by a remote server. it also uses the basic 'Write'
// method (of the io.Writer interface) to send simple text messages to the server,
// which the server echoes back.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var server = flag.String("server", "localhost", "name/address of echo server.")
var port = flag.Int("port", 5432, "server port number to connect to.")

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
			case msg := <-svc.Received():
				fmt.Println("received:", msg)
			case <-svc.Closed():
				fmt.Println("service closed.")
				return
			}
		}
	}()

	userInput := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("echocli: ")
		userInput.Scan()
		if in := userInput.Text(); in == "quit" {
			break
		}

		svc.Write([]byte(in))
	}
}
