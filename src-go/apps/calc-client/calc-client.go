package main

import "github.com/richardmillen/etude-2-net-patterns/src-go/core"

func main() {
	d := core.NewDialer("tcp", "localhost:5432")
	defer d.Close()

}
