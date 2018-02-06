package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var host = flag.String("server", "localhost", "name of log-server hsot.")
var port = flag.Int("port", 5959, "port number to connect to.")
var count = flag.Int("count", 100, "number of times run test function.")
var pause = flag.Bool("pause", false, "whether to pause briefly between test function calls.")

var (
	alerts    = 1
	criticals = 1
	errors    = 1
	warnings  = 1
	notices   = 1
)

func init() {
	log.SetPrefix(fmt.Sprintf("log-client [%d]: ", os.Getpid()))
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("starting log/trace client (addr: %s)...\n", addr)

	d := netx.NewDialer("tcp", addr)
	defer d.Close()

	pub := netx.NewPublisher(d)
	defer pub.Close()

	func() {
		mainLog := netx.Start(pub, "main")
		defer mainLog.Close()

		for n := 1; n <= *count; n++ {
			msg := fmt.Sprintf("running test #%d...", n)

			log.Println(msg)
			mainLog.Print(netx.Notice, msg)

			foo(mainLog)
		}

		mainLog.Print(netx.Notice, "finished running tests.")
	}()

	log.Println("stats:")
	log.Printf("\talerts: %d\n", alerts)
	log.Printf("\tcriticals: %d\n", criticals)
	log.Printf("\terrors: %d\n", errors)
	log.Printf("\twarnings: %d\n", warnings)
	log.Printf("\tnotices: %d\n", notices)
}

func foo(parentLog *netx.Logger) {
	fooLog := netx.StartChild(parentLog, "foo")
	defer fooLog.Close()

	if *pause {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}

	switch rand.Intn(10) {
	case 0:
		fooLog.Printf(netx.Alert, "alert message %d", alerts)
		alerts++
	case 1, 2:
		fooLog.Printf(netx.Critical, "critical message %d", criticals)
		criticals++
	case 3, 4, 5:
		baz(fooLog)
	default:
		bar(fooLog)
	}
}

func bar(parentLog *netx.Logger) {
	barLog := netx.StartChild(parentLog, "bar")
	defer barLog.Close()

	if *pause {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}

	switch rand.Intn(10) {
	case 0:
		barLog.Printf(netx.Error, "error message %d", errors)
		errors++
	case 1:
		barLog.Printf(netx.Warning, "warning message %d", warnings)
		warnings++
	case 2, 3, 4:
		barLog.Print(netx.Debug, "debug message")
	default:
		baz(barLog)
	}
}

func baz(parentLog *netx.Logger) {
	bazLog := netx.StartChild(parentLog, "baz")
	defer bazLog.Close()

	if *pause {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}

	switch rand.Intn(10) {
	case 0, 1, 2, 3, 4:
		bazLog.Printf(netx.Notice, "notice message %d", notices)
		notices++
	case 5, 6:
		bazLog.Print(netx.Debug, "debug message")
	default:
		bazLog.Print(netx.Info, "info message")
	}
}
