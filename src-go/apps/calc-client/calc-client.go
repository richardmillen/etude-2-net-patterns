package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"

	"github.com/richardmillen/etude-2-net-patterns/src-go/core"
)

const (
	cmdHelp string = "help"
	cmdQuit string = "quit"
)

var scanner = bufio.NewScanner(os.Stdin)

func main() {
	d := core.NewDialer("tcp", "localhost:5432")
	defer d.Close()

	s := core.NewService(d, &calcClient{})
	defer s.Close()

	s.Connect(func(q *core.Queue) error {
		return core.ErrNoImpl
	})
	s.Error(func(error) error {
		return core.ErrNoImpl
	})
	s.Recv(func(v interface{}) error {
		return core.ErrNoImpl
	})
	s.Start()

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
}
