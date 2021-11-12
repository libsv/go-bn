package zmq

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-zeromq/zmq4"
)

// MessageFunc a func in which the message is passed to.
type MessageFunc func([][]byte)

// ErrorFunc a func in which an error is passed to.
type ErrorFunc func(err error)

// Topic a subscription topic.
type Topic string

// Subscription topics.
const (
	TopicHashTx                  Topic = "hashtx"
	TopicHashBlock               Topic = "hashblock"
	TopicInvalidTx               Topic = "invalidtx"
	TopicDicardFromMempool       Topic = "discardfrommempool"
	TopicRemovedFromMempoolBlock Topic = "removedfrommempoolblock"

	TopicRawTx    Topic = "rawtx"
	TopicRawBlock Topic = "rawblock"
)

type zmq struct {
	mu            sync.RWMutex
	conn          zmq4.Socket
	connected     bool
	onErrFn       ErrorFunc
	cfg           *nodeMqCfg
	subscriptions map[Topic]MessageFunc
}

// NodeMQ interfaces connecting and subscribing to a bitcoin node NodeMQ connection.
type NodeMQ interface {
	Connect(ctx context.Context) error
	Subscribe(topic Topic, fn MessageFunc) error
	Unsubscribe(topic Topic) error
}

// NewNodeMQ build and return a new zmq.ZMQ configured via the provided opt funcs.
func NewNodeMQ(oo ...NodeMQOptFunc) NodeMQ {
	cfg := &nodeMqCfg{
		optionValue: "hash",
		errorFn:     defaultOnError,
		topics: map[Topic]bool{
			TopicHashBlock:               true,
			TopicHashTx:                  true,
			TopicDicardFromMempool:       true,
			TopicRemovedFromMempoolBlock: true,
			TopicInvalidTx:               true,

			TopicRawTx:    false,
			TopicRawBlock: false,
		},
	}
	for _, o := range oo {
		o(cfg)
	}

	return &zmq{
		cfg:           cfg,
		subscriptions: make(map[Topic]MessageFunc),
		onErrFn:       cfg.errorFn,
	}
}

// Connect to the bitcoin node 0MQ.
func (z *zmq) Connect(ctx context.Context) error {
	if err := z.cfg.validate(); err != nil {
		return err
	}

	z.conn = zmq4.NewSub(ctx, zmq4.WithID(zmq4.SocketIdentity("sub")))
	if err := z.conn.Dial(z.cfg.host); err != nil {
		return err
	}

	defer func() {
		if !z.connected {
			return
		}

		if err := z.conn.Close(); err != nil {
			z.onErrFn(err)
		}
		z.connected = false
	}()

	if err := z.conn.SetOption(zmq4.OptionSubscribe, z.cfg.optionValue); err != nil {
		return err
	}

	if z.cfg.raw {
		if err := z.conn.SetOption(zmq4.OptionSubscribe, "raw"); err != nil {
			return err
		}
	}

	for {
		msg, err := z.conn.Recv()
		if err != nil {
			z.onErrFn(err)
		}
		if !z.connected {
			z.connected = true
		}
		func() {
			z.mu.RLock()
			defer z.mu.RUnlock()

			fn, ok := z.subscriptions[Topic(msg.Frames[0])]
			if ok {
				go fn(msg.Frames)
			}
		}()
	}
}

// Subscribe to a topic on a bitcoin node 0MQ.
func (z *zmq) Subscribe(topic Topic, fn MessageFunc) error {
	if ok := z.cfg.topics[topic]; !ok {
		return fmt.Errorf("unrecognised topic: %s", topic)
	}

	if !z.cfg.allowOverwrite {
		if _, ok := z.subscriptions[topic]; ok {
			return fmt.Errorf("already subscribed to %s", topic)
		}
	}

	z.mu.Lock()
	defer z.mu.Unlock()

	z.subscriptions[topic] = fn
	return nil
}

// Unsubscribe from a topic on the bitcoin ndoe 0MQ.
func (z *zmq) Unsubscribe(topic Topic) error {
	if ok := z.cfg.topics[topic]; !ok {
		return fmt.Errorf("unrecognised topic: %s", topic)
	}

	z.mu.Lock()
	defer z.mu.Unlock()

	delete(z.subscriptions, topic)
	return nil
}

func defaultOnError(err error) {
	fmt.Fprintln(os.Stderr, err)
}
