package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var port = flag.Int("port", 5959, "Port number to listen at.")
var severity = flag.String("severity", netx.DebugTopic, "the severity topic (HACK: countdown from minimum severity value).")

func init() {
	log.SetPrefix("log-collector: ")
}

func main() {
	flag.Parse()

	log.Printf("starting log/trace collector (port: %d)...\n", *port)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	listener, err := netx.ListenTCP("tcp", addr)
	check.Error(err)
	defer listener.Close()

	sub := netx.NewSubscriber(listener, *severity)
	defer sub.Close()

	sub.Error(func(err error) error {
		log.Println("collector error:", err)
		return nil
	})

	sub.Subscribe(func(m *netx.Message) (err error) {
		log.Printf("%s: %s\n", m.Topic, string(m.Body))
		return
	})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	<-sigint
	log.Println("server interrupted.")
}
