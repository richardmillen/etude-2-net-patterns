package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

var port = flag.Int("port", 5432, "port number to listen at")
var sendErrors = flag.Bool("send-errors", true, "specify whether the server should return an error in response to an invalid message")

func main() {
	flag.Parse()

	log.Println("configuring server states...")

	recvState := &fsm.State{
		Name: "receiving",
		Accepts: []fsm.Input{
			&fsm.String{
				Hint:  "accept 'hello'",
				Match: "hello",
			},
			&fsm.String{
				Hint:  "accept 'hi'",
				Match: "hi",
			},
		},
	}
	errorState := &fsm.State{
		Name: "invalid message",
		Accepts: []fsm.Input{
			&fsm.Any{},
		},
		Substates: []fsm.State{
			recvState,
		},
	}

	log.Println("starting service...")

	svc, err := service.ListenTCP("tcp", fmt.Sprintf(":%d", *port))
	check.Error(err)

	svc.Initial(recvState)

	//sigint := make(chan os.Signal, 1)
	//signal.Notify(sigint, os.Interrupt)

	log.Println("server listening...")

	for {
		select {
		case r := <-svc.Receiver():
			go func(r *service.Receiver) {
				switch in := r.Input.(type) {
				case *fsm.String:
					log.Println("received:", in)
					r.Conn.Write([]byte("world"))
				}
			}(r)
		case <-svc.Interrupt():
			log.Println("server interrupted.")
			return
		}
	}

}
