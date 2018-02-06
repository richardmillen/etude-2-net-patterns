package main

import "github.com/richardmillen/etude-2-net-patterns/src-go/core"

type calcServer struct {
}

func (c *calcServer) Greet(q *core.Queue) error {
	return core.ErrNoImpl
}

func (c *calcServer) Send(q *core.Queue, v interface{}) error {
	return core.ErrNoImpl
}

func (c *calcServer) Recv(q *core.Queue) (interface{}, error) {
	return nil, core.ErrNoImpl
}
