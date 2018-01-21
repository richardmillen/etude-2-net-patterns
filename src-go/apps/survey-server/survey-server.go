package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	source := ":1202"
	conn, err := net.ListenPacket("udp", source)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		fmt.Println("waiting...")

		buf := make([]byte, 8)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			continue
		}

		fmt.Printf("echo: %v\n", string(buf))
		_, err = conn.WriteTo(buf, addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		}
	}
}
