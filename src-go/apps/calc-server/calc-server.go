package main

import (
	"net"

	"github.com/richardmillen/etude-2-net-patterns/src-go/check"
	"github.com/richardmillen/etude-2-net-patterns/src-go/core"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":5432")
	check.Error(err)

	listener, err := core.ListenTCP("tcp", addr)
	check.Error(err)
	defer listener.Close()

}
