package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var port = flag.Int("port", 5959, "Port number to listen at.")
var severity = flag.String("severity", "debug", "the severity level.")

func init() {
	log.SetPrefix("log-collector: ")
}

func main() {
	flag.Parse()

	// ...

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}
