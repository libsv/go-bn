package bn

import "time"

// BitcoinClientOptFunc for setting bitcoin client options.
type BitcoinClientOptFunc func(c *clientOpts)

type clientOpts struct {
	timeout   time.Duration
	host      string
	username  string
	password  string
	cache     bool
	isMainnet bool
}

// WithTimeout set the timeout for the http client.
func WithTimeout(seconds time.Duration) BitcoinClientOptFunc {
	return func(c *clientOpts) {
		c.timeout = seconds
	}
}

// WithCache enable response caching.
func WithCache() BitcoinClientOptFunc {
	return func(c *clientOpts) {
		c.cache = true
	}
}

// WithHost set the bitcoin node host.
func WithHost(host string) BitcoinClientOptFunc {
	return func(c *clientOpts) {
		c.host = host
	}
}

// WithCreds set the bitcoin node credentials.
func WithCreds(username, password string) BitcoinClientOptFunc {
	return func(c *clientOpts) {
		c.username = username
		c.password = password
	}
}

// WithMainnet set whether or not the node is a mainnet node.
func WithMainnet() BitcoinClientOptFunc {
	return func(c *clientOpts) {
		c.isMainnet = true
	}
}
