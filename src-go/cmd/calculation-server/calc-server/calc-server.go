package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var port = flag.Int("port", 5432, "port number to listen at")

func main() {
	log.Println("starting calc-server...")

	// ....

	log.Println("calc-server started.")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}
