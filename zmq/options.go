package zmq

import "errors"

type nodeMqCfg struct {
	host           string
	raw            bool
	topics         map[Topic]bool
	optionValue    string
	allowOverwrite bool
	errorFn        ErrorFunc
}

func (c *nodeMqCfg) validate() error {
	if c.host == "" {
		return errors.New("host cannot be empty")
	}

	return nil
}

// NodeMQOptFunc option func.
type NodeMQOptFunc func(o *nodeMqCfg)

// WithRaw listen and allow subscribing to `rawtx` and `rawblock` messages.
func WithRaw() NodeMQOptFunc {
	return func(o *nodeMqCfg) {
		o.topics[TopicRawBlock] = true
		o.topics[TopicRawTx] = true
		o.raw = true
	}
}

// WithSubscribeOptionValue set the option value.
func WithSubscribeOptionValue(ov string) NodeMQOptFunc {
	return func(o *nodeMqCfg) {
		o.optionValue = ov
	}
}

// WithHost set the host to connect to. Expects the following format:
// tcp://hostname:port
func WithHost(host string) NodeMQOptFunc {
	return func(o *nodeMqCfg) {
		o.host = host
	}
}

// WithSubscriptionOverwrite allows the reassignment of a topic subscribe, without first
// unsubscribing.
//
// When set, the following no longer becomes an error scenario:
//   z := zmq.NewZMQ(zmq.WithHost(...), zmq.WithSubscriptionOverwrite())
//   if err := z.Subscribe(zmq.TopicHashTx, func([][]byte){}); err != nil {}
//   if err := z.Subscribe(zmq.TopicHashTx, func([][]byte){}); err != nil {}
func WithSubscriptionOverwrite() NodeMQOptFunc {
	return func(o *nodeMqCfg) {
		o.allowOverwrite = true
	}
}

// WithErrorHandler sets an error handler func.
func WithErrorHandler(fn ErrorFunc) NodeMQOptFunc {
	return func(o *nodeMqCfg) {
		o.errorFn = fn
	}
}
