package zmq

import "errors"

// Standard errors.
var (
	ErrInvalidTopic      = errors.New("invalid topic")
	ErrAlreadySubscribed = errors.New("already subscribed to topic")
	ErrHostEmpty         = errors.New("host cannot be empty")
)
