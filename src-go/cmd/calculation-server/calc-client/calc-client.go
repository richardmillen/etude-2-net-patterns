package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

const (
	cmdHelp string = "help"
	cmdQuit string = "quit"
)

var file = flag.String("file", "", "specify a file that contains one or more calculations to be sent to the server e.g. calcs.txt")
var server = flag.String("server", "localhost", "the name/address of the calc-server host")
var port = flag.Int("port", 5432, "the port to connect to on the host")

var scanner = bufio.NewScanner(os.Stdin)

func main() {
	// ...

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
			// ...
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
