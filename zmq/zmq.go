package zmq

import (
	"context"
	"sync"

	"github.com/go-zeromq/zmq4"
)

type zmq struct {
	mu        sync.RWMutex
	address   string
	s         zmq4.Socket
	connected bool
	opts      *zmqOpts
}

type ZMQ interface {
	Connect(ctx context.Context) error
	Subscribe() error
	Unsubscribe() error
}

func NewZMQ(oo ...ZMQOptFunc) ZMQ {
	opts := &zmqOpts{
		additionalTopics: make([]string, 0),
		optionValue:      "hash",
	}
	for _, o := range oo {
		o(opts)
	}
	return &zmq{
		opts: opts,
	}
}

func (z *zmq) Connect(ctx context.Context) error {
	return nil
}

func (z *zmq) Subscribe() error {
	panic("not implemented") // TODO: Implement
}

func (z *zmq) Unsubscribe() error {
	panic("not implemented") // TODO: Implement
}
