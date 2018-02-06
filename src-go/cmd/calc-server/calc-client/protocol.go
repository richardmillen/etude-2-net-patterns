package main

import "github.com/richardmillen/etude-2-net-patterns/src-go/core"

type calcClient struct {
}

func (c *calcClient) Greet(q *core.Queue) error {
	return core.ErrNoImpl
}

func (c *calcClient) Send(q *core.Queue, v interface{}) error {
	return core.ErrNoImpl
}

func (c *calcClient) Recv(q *core.Queue) (interface{}, error) {
	return nil, core.ErrNoImpl
}
