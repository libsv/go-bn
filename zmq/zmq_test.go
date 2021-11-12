package zmq_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/go-zeromq/zmq4"
	"github.com/libsv/go-bn/mocks"
	"github.com/libsv/go-bn/zmq"
	"github.com/stretchr/testify/assert"
)

func TestNodeMQ_Subscribe(t *testing.T) {
	tests := map[string]struct {
		topics []zmq.Topic
		opts   []zmq.NodeMQOptFunc
		expErr error
	}{
		"successful subscription hashtx": {
			topics: []zmq.Topic{zmq.TopicHashTx},
		},
		"successful subscription hashblock": {
			topics: []zmq.Topic{zmq.TopicHashBlock},
		},
		"successful subscription discardfrommempool": {
			topics: []zmq.Topic{zmq.TopicDiscardFromMempool},
		},
		"successful subscription invalidtx": {
			topics: []zmq.Topic{zmq.TopicInvalidTx},
		},
		"successful subscription removedfrommempoolblock": {
			topics: []zmq.Topic{zmq.TopicRemovedFromMempoolBlock},
		},
		"successful subscription to all non-raw topics": {
			topics: []zmq.Topic{
				zmq.TopicHashTx,
				zmq.TopicHashBlock,
				zmq.TopicInvalidTx,
				zmq.TopicDiscardFromMempool,
				zmq.TopicRemovedFromMempoolBlock,
			},
		},
		"error subscribing to rawtx without being enabled": {
			topics: []zmq.Topic{zmq.TopicRawTx},
			expErr: errors.New("invalid topic: rawtx"),
		},
		"error subscribing to rawblock without being enabled": {
			topics: []zmq.Topic{zmq.TopicRawBlock},
			expErr: errors.New("invalid topic: rawblock"),
		},
		"error subscribing to all topics without raw enabled": {
			topics: []zmq.Topic{
				zmq.TopicHashTx,
				zmq.TopicHashBlock,
				zmq.TopicInvalidTx,
				zmq.TopicDiscardFromMempool,
				zmq.TopicRemovedFromMempoolBlock,
				zmq.TopicRawTx,
				zmq.TopicRawBlock,
			},
			expErr: errors.New("invalid topic: rawtx"),
		},
		"successful subscription rawtx when enabled": {
			topics: []zmq.Topic{zmq.TopicRawTx},
			opts:   []zmq.NodeMQOptFunc{zmq.WithRaw()},
		},
		"successful subscription rawblock when enabled": {
			topics: []zmq.Topic{zmq.TopicRawBlock},
			opts:   []zmq.NodeMQOptFunc{zmq.WithRaw()},
		},
		"successful subscription to all topics with raw enabled": {
			topics: []zmq.Topic{
				zmq.TopicHashTx,
				zmq.TopicHashBlock,
				zmq.TopicInvalidTx,
				zmq.TopicDiscardFromMempool,
				zmq.TopicRemovedFromMempoolBlock,
				zmq.TopicRawTx,
				zmq.TopicRawBlock,
			},
			opts: []zmq.NodeMQOptFunc{zmq.WithRaw()},
		},
		"error resubscribing to topic without enabling": {
			topics: []zmq.Topic{
				zmq.TopicHashTx,
				zmq.TopicHashTx,
			},
			expErr: errors.New("already subscribed to topic: hashtx"),
		},
		"successful resubscription to topic when enabled": {
			topics: []zmq.Topic{
				zmq.TopicHashTx,
				zmq.TopicHashTx,
			},
			opts: []zmq.NodeMQOptFunc{zmq.WithSubscriptionOverwrite()},
		},
		"error subscribing to topic that does not exist": {
			topics: []zmq.Topic{
				zmq.Topic("oh hello there"),
			},
			expErr: errors.New("invalid topic: oh hello there"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			z := zmq.NewNodeMQ(append(test.opts, zmq.WithCustomZMQSocket(&mocks.SocketMock{}))...)
			err := func() error {
				for _, topic := range test.topics {
					if err := z.Subscribe(topic, func([][]byte) {}); err != nil {
						return err
					}
				}
				return nil
			}()
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNodeMQ_Unsubscribe(t *testing.T) {
	tests := map[string]struct {
		subscribedTo    []zmq.Topic
		unsubscribeFrom []zmq.Topic
		opts            []zmq.NodeMQOptFunc
		expErr          error
	}{
		"successful unsubscribe": {
			subscribedTo:    []zmq.Topic{zmq.TopicHashTx},
			unsubscribeFrom: []zmq.Topic{zmq.TopicHashTx},
		},
		"no error when unsubscribing from a topic no subscribed to": {
			subscribedTo:    []zmq.Topic{zmq.TopicHashTx},
			unsubscribeFrom: []zmq.Topic{zmq.TopicInvalidTx},
		},
		"error when unsubscribing from raw topic when not enabled": {
			subscribedTo:    []zmq.Topic{zmq.TopicHashTx},
			unsubscribeFrom: []zmq.Topic{zmq.TopicRawTx},
			expErr:          errors.New("invalid topic: rawtx"),
		},
		"successful unsubscribe from raw topic when enabled": {
			subscribedTo:    []zmq.Topic{zmq.TopicRawTx},
			unsubscribeFrom: []zmq.Topic{zmq.TopicRawTx},
			opts:            []zmq.NodeMQOptFunc{zmq.WithRaw()},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			z := zmq.NewNodeMQ(append(test.opts, zmq.WithCustomZMQSocket(&mocks.SocketMock{}))...)
			for _, topic := range test.subscribedTo {
				assert.NoError(t, z.Subscribe(topic, func([][]byte) {}))
			}

			assert.True(t, len(test.unsubscribeFrom) > 0, "test %s has not declare a topic to unsub from")

			err := func() error {
				for _, topic := range test.unsubscribeFrom {
					if err := z.Unsubscribe(topic); err != nil {
						return err
					}
				}
				return nil
			}()
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNodeMQ_Connect(t *testing.T) {
	t.Parallel()
	type option struct {
		name  string
		value interface{}
	}

	type message struct {
		msg zmq4.Msg
		err error
	}

	tests := map[string]struct {
		topics                []zmq.Topic
		host                  string
		messages              []message
		opts                  []zmq.NodeMQOptFunc
		socketDialFn          func(string) error
		setOptionFunc         func(string, interface{}) error
		closeFunc             func() error
		recvFunc              func() (zmq4.Msg, error)
		expCounts             map[zmq.Topic]int
		expOptions            []option
		expConnectError       error
		expErrorHandlerErrors []error
	}{
		"error invalid host name": {
			host: "",
			closeFunc: func() error {
				return nil
			},
			expCounts:       map[zmq.Topic]int{},
			expConnectError: errors.New("host cannot be empty"),
		},
		"error dailing socket is returned": {
			host: "tcp://localhost:12345",
			socketDialFn: func(s string) error {
				return errors.New("YIKES")
			},
			expCounts:       map[zmq.Topic]int{},
			expConnectError: errors.New("YIKES"),
		},
		"error setting option is returned": {
			host: "tcp://localhost:12345",
			socketDialFn: func(s string) error {
				return nil
			},
			setOptionFunc: func(string, interface{}) error {
				return errors.New("no options 4 u")
			},
			expOptions: []option{{
				name:  "SUBSCRIBE",
				value: "hash",
			}},
			expCounts:       map[zmq.Topic]int{},
			expConnectError: errors.New("no options 4 u"),
		},
		"error on close is reported": {
			host: "tcp://localhost:12345",
			messages: []message{{
				msg: zmq4.Msg{
					Frames: func() [][]byte {
						return [][]byte{
							[]byte("hashtx"),
						}
					}(),
				},
			}},
			socketDialFn: func(s string) error {
				return nil
			},
			setOptionFunc: func(string, interface{}) error {
				return nil
			},
			recvFunc: func() (zmq4.Msg, error) {
				return zmq4.Msg{}, nil
			},
			closeFunc: func() error {
				return errors.New("oh no")
			},
			expOptions: []option{{
				name:  "SUBSCRIBE",
				value: "hash",
			}},
			expCounts:             map[zmq.Topic]int{},
			expErrorHandlerErrors: []error{errors.New("oh no")},
		},
		"error with received messages are reported": {
			host: "tcp://localhost:12345",
			messages: []message{{
				err: errors.New("first error"),
			}, {
				err: errors.New("second error"),
			}, {
				msg: zmq4.Msg{
					Frames: func() [][]byte {
						return [][]byte{
							[]byte("hashtx"),
						}
					}(),
				},
			}, {
				err: errors.New("third error"),
			}},
			socketDialFn: func(s string) error {
				return nil
			},
			setOptionFunc: func(string, interface{}) error {
				return nil
			},
			recvFunc: func() (zmq4.Msg, error) {
				return zmq4.Msg{}, nil
			},
			closeFunc: func() error {
				return nil
			},
			expOptions: []option{{
				name:  "SUBSCRIBE",
				value: "hash",
			}},
			expCounts: map[zmq.Topic]int{},
			expErrorHandlerErrors: []error{
				errors.New("first error"),
				errors.New("second error"),
				errors.New("third error"),
			},
		},
		"raw option present when WithRaw is called": {
			host: "tcp://localhost:12345",
			opts: []zmq.NodeMQOptFunc{zmq.WithRaw()},
			socketDialFn: func(s string) error {
				return nil
			},
			setOptionFunc: func(string, interface{}) error {
				return nil
			},
			recvFunc: func() (zmq4.Msg, error) {
				return zmq4.Msg{}, nil
			},
			closeFunc: func() error {
				return nil
			},
			expOptions: []option{{
				name:  "SUBSCRIBE",
				value: "hash",
			}, {
				name:  "SUBSCRIBE",
				value: "raw",
			}},
			expCounts: map[zmq.Topic]int{},
		},
		"only messages on subscribed topics are relayed": {
			host: "tcp://localhost:12345",
			topics: []zmq.Topic{
				zmq.TopicHashTx,
				zmq.TopicInvalidTx,
			},
			messages: []message{{
				msg: zmq4.Msg{
					Frames: func() [][]byte {
						return [][]byte{
							[]byte("hashtx"),
						}
					}(),
				},
			}, {
				msg: zmq4.Msg{
					Frames: func() [][]byte {
						return [][]byte{
							[]byte("hashtx"),
						}
					}(),
				},
			}, {
				msg: zmq4.Msg{
					Frames: func() [][]byte {
						return [][]byte{
							[]byte("discardfromempool"),
						}
					}(),
				},
			}, {
				msg: zmq4.Msg{
					Frames: func() [][]byte {
						return [][]byte{
							[]byte("removedfrommempoolblock"),
						}
					}(),
				},
			}, {
				msg: zmq4.Msg{
					Frames: func() [][]byte {
						return [][]byte{
							[]byte("invalidtx"),
						}
					}(),
				},
			}},
			socketDialFn: func(s string) error {
				return nil
			},
			setOptionFunc: func(string, interface{}) error {
				return nil
			},
			recvFunc: func() (zmq4.Msg, error) {
				return zmq4.Msg{}, nil
			},
			closeFunc: func() error {
				return nil
			},
			expOptions: []option{{
				name:  "SUBSCRIBE",
				value: "hash",
			}},
			expCounts: map[zmq.Topic]int{
				zmq.TopicHashTx:    2,
				zmq.TopicInvalidTx: 1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var options []option
			socket := &mocks.SocketMock{
				DialFunc: func(addr string) error {
					assert.Equal(t, test.host, addr)
					return test.socketDialFn(addr)
				},
				SetOptionFunc: func(opt string, v interface{}) error {
					options = append(options, option{opt, v})
					return test.setOptionFunc(opt, v)
				},
				RecvFunc: func() (zmq4.Msg, error) {
					if len(test.messages) == 0 {
						return zmq4.Msg{}, context.Canceled
					}
					defer func() { test.messages = test.messages[1:] }()

					if _, err := test.recvFunc(); err != nil {
						return zmq4.Msg{}, err
					}
					return test.messages[0].msg, test.messages[0].err
				},
				CloseFunc: test.closeFunc,
			}

			var errHandlerErrs []error

			z := zmq.NewNodeMQ(
				append(
					test.opts,
					zmq.WithCustomZMQSocket(socket),
					zmq.WithErrorHandler(func(err error) {
						errHandlerErrs = append(errHandlerErrs, err)
					}),
					zmq.WithHost(test.host),
				)...,
			)

			var total int
			for _, v := range test.expCounts {
				total += v
			}

			m := map[zmq.Topic]int{}
			var mu sync.Mutex
			var wg sync.WaitGroup
			wg.Add(total)
			for _, topic := range test.topics {
				topic := topic
				assert.NoError(t, z.Subscribe(topic, func(msg [][]byte) {
					defer wg.Done()
					mu.Lock()
					defer mu.Unlock()
					m[topic] = m[topic] + 1
				}))
			}

			err := z.Connect()
			if test.expConnectError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expConnectError.Error())
			} else {
				assert.NoError(t, err)
			}
			wg.Wait()

			assert.Equal(t, test.expOptions, options)
			assert.Equal(t, test.expCounts, m)
			assert.Equal(t, test.expErrorHandlerErrors, errHandlerErrs)
		})
	}
}
