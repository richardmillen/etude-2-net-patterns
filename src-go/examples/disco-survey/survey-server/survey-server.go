//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

// default port is 0 (zero) in order to use ephemeral port.
var port = flag.Int("echo-port", 0, "port number to listen for echo requests at.")

func init() {
	log.SetPrefix("survey-server: ")
}

func main() {
	flag.Parse()

	// ....

	log.Println("survey server started.")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}

func echoPortDesc() string {
	if *port == 0 {
		return "ephemeral port"
	}
	return fmt.Sprintf("port %d", *port)
}
