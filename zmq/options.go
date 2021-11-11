package zmq

type zmqOpts struct {
	host             string
	raw              bool
	additionalTopics []string
	optionValue      string
}

type ZMQOptFunc func(o *zmqOpts)

func WithRaw() ZMQOptFunc {
	return func(o *zmqOpts) {
		o.additionalTopics = append(o.additionalTopics, "rawblock", "rawtx")
		o.raw = true
	}
}

func WithSubscribeOptionValue(ov string) ZMQOptFunc {
	return func(o *zmqOpts) {
		o.optionValue = ov
	}
}

// WithHost set the host to connect to. Expects the following format:
// tcp://hostname:port
func WithHost(host string) ZMQOptFunc {
	return func(o *zmqOpts) {
		o.host = host
	}
}
