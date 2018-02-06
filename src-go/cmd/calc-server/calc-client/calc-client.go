package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
)

const (
	cmdHelp string = "help"
	cmdQuit string = "quit"
)

var server = flag.String("server", "localhost", "the name/address of the calc-server host")
var port = flag.Int("port", 5432, "the port to connect to on the host")

var scanner = bufio.NewScanner(os.Stdin)

func main() {
	log.Println("starting calc-client...")

	d := netx.NewDialer("tcp", fmt.Sprintf("%s:%d", *server, *port))
	defer d.Close()

	s := netx.NewService(d, &calcClient{})
	defer s.Close()

	log.Println("calc-client started.")

	printHelp()

	for {
		input := getInput()
		switch input {
		case cmdHelp:
			printHelp()
		case cmdQuit:
			fmt.Println("quitting...")
			return
		default:
			_, err := s.Write([]byte(input))
			check.Error(err)
		}
	}
}

func getInput() string {
	fmt.Printf("calc: ")
	scanner.Scan()
	return scanner.Text()
}

func printHelp() {
	fmt.Println()
	fmt.Println("----------------------------------------------------------------------------------------------")
	fmt.Println("commands:")
	fmt.Printf("\t%s\tprint this help message.\n", cmdHelp)
	fmt.Printf("\t%s\texit client.\n", cmdQuit)
	fmt.Println()
	fmt.Println("all other input will be sent to the server which will calculate the result or return an error.")
	fmt.Println("\te.g.")
	fmt.Println("\t2/2")
	fmt.Println("\tresult: 1")
	fmt.Println("\t2*2")
	fmt.Println("\tresult: 4")
	fmt.Println("\t2+2")
	fmt.Println("\tresult: 4")
	fmt.Println("\t2-2")
	fmt.Println("\tresult: 0")
	fmt.Println("----------------------------------------------------------------------------------------------")
	fmt.Println()
}
