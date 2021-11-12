package zmq

import (
	"context"
	"errors"
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
	TopicDiscardFromMempool      Topic = "discardfrommempool"
	TopicRemovedFromMempoolBlock Topic = "removedfrommempoolblock"

	TopicRawTx    Topic = "rawtx"
	TopicRawBlock Topic = "rawblock"
)

type nodeMq struct {
	mu            sync.RWMutex
	conn          zmq4.Socket
	connected     bool
	onErrFn       ErrorFunc
	cfg           *nodeMqCfg
	subscriptions map[Topic]MessageFunc
}

// NodeMQ interfaces connecting and subscribing to a bitcoin node NodeMQ connection.
type NodeMQ interface {
	Connect() error
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
			TopicDiscardFromMempool:      true,
			TopicRemovedFromMempoolBlock: true,
			TopicInvalidTx:               true,

			TopicRawTx:    false,
			TopicRawBlock: false,
		},
	}
	for _, o := range oo {
		o(cfg)
	}

	if cfg.zmqSocket == nil {
		cfg.zmqSocket = zmq4.NewSub(cfg.ctx, zmq4.WithID(zmq4.SocketIdentity("sub")))
	}

	return &nodeMq{
		cfg:           cfg,
		subscriptions: make(map[Topic]MessageFunc),
		onErrFn:       cfg.errorFn,
		conn:          cfg.zmqSocket,
	}
}

// Connect to the bitcoin node 0MQ.
func (n *nodeMq) Connect() error {
	if err := n.cfg.validate(); err != nil {
		return err
	}

	if err := n.conn.Dial(n.cfg.host); err != nil {
		return err
	}

	defer func() {
		if !n.connected {
			return
		}

		if err := n.conn.Close(); err != nil {
			n.onErrFn(err)
		}
		n.connected = false
	}()

	if err := n.conn.SetOption(zmq4.OptionSubscribe, n.cfg.optionValue); err != nil {
		return err
	}

	if n.cfg.raw {
		if err := n.conn.SetOption(zmq4.OptionSubscribe, "raw"); err != nil {
			return err
		}
	}

	for {
		msg, err := n.conn.Recv()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			n.onErrFn(err)
			continue
		}
		n.connected = true
		func() {
			n.mu.RLock()
			defer n.mu.RUnlock()

			fn, ok := n.subscriptions[Topic(msg.Frames[0])]
			if ok {
				go fn(msg.Frames)
			}
		}()
	}
}

// Subscribe to a topic on a bitcoin node 0MQ.
func (n *nodeMq) Subscribe(topic Topic, fn MessageFunc) error {
	if ok := n.cfg.topics[topic]; !ok {
		return fmt.Errorf("%w: %s", ErrInvalidTopic, topic)
	}

	if !n.cfg.allowOverwrite {
		if _, ok := n.subscriptions[topic]; ok {
			return fmt.Errorf("%w: %s", ErrAlreadySubscribed, topic)
		}
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.subscriptions[topic] = fn
	return nil
}

// Unsubscribe from a topic on the bitcoin ndoe 0MQ.
func (n *nodeMq) Unsubscribe(topic Topic) error {
	if ok := n.cfg.topics[topic]; !ok {
		return fmt.Errorf("%w: %s", ErrInvalidTopic, topic)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.subscriptions, topic)
	return nil
}

func defaultOnError(err error) {
	fmt.Fprintln(os.Stderr, err)
}
