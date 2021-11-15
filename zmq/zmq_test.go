package zmq_test

import (
	"context"
	"encoding/hex"
	"errors"
	"sync"
	"testing"

	"github.com/go-zeromq/zmq4"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-bn/mocks"
	"github.com/libsv/go-bn/zmq"
	"github.com/libsv/go-bt/v2"
	"github.com/stretchr/testify/assert"
)

func TestNodeMQ_Subscribe(t *testing.T) {
	t.Parallel()

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
					if err := z.Subscribe(topic, func(context.Context, [][]byte) {}); err != nil {
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
	t.Parallel()

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
				assert.NoError(t, z.Subscribe(topic, func(context.Context, [][]byte) {}))
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
					zmq.WithErrorHandler(func(_ context.Context, err error) {
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
				assert.NoError(t, z.Subscribe(topic, func(_ context.Context, msg [][]byte) {
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

func TestNodeMQ_SubscribeHashTx(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		messages  []zmq4.Msg
		expHashes []string
	}{
		"successful messages": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashTx)
					body, err := hex.DecodeString("b9037b78403cc4e76a06060b4e9d1e1cdf7c85ce7cb6f074e1f9ed2fb6aa10a6")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashTx)
					body, err := hex.DecodeString("4f69dda2c8b2eb01ebf002d559e15e7ac183acf51730e2bf889d4638864675a9")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expHashes: []string{
				"4f69dda2c8b2eb01ebf002d559e15e7ac183acf51730e2bf889d4638864675a9",
				"b9037b78403cc4e76a06060b4e9d1e1cdf7c85ce7cb6f074e1f9ed2fb6aa10a6",
			},
		},
		"irrelevant messages are ignored": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashTx)
					body, err := hex.DecodeString("b9037b78403cc4e76a06060b4e9d1e1cdf7c85ce7cb6f074e1f9ed2fb6aa10a6")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("4f69dda2c8b2eb01ebf002d559e15e7ac183acf51730e2bf889d4638864675a9")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashTx)
					body, err := hex.DecodeString("4f69dda2c8b2eb01ebf002d559e15e7ac183acf51730e2bf889d4638864675a9")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expHashes: []string{
				"4f69dda2c8b2eb01ebf002d559e15e7ac183acf51730e2bf889d4638864675a9",
				"b9037b78403cc4e76a06060b4e9d1e1cdf7c85ce7cb6f074e1f9ed2fb6aa10a6",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var errHandlerErrs []error
			socket := &mocks.SocketMock{
				DialFunc: func(addr string) error {
					return nil
				},
				SetOptionFunc: func(opt string, v interface{}) error {
					return nil
				},
				RecvFunc: func() (zmq4.Msg, error) {
					if len(test.messages) == 0 {
						return zmq4.Msg{}, context.Canceled
					}
					defer func() { test.messages = test.messages[1:] }()

					return test.messages[0], nil
				},
				CloseFunc: func() error {
					return nil
				},
			}
			c := zmq.NewNodeMQ(
				zmq.WithHost("tcp://localhost:12345"),
				zmq.WithCustomZMQSocket(socket),
				zmq.WithErrorHandler(func(ctx context.Context, err error) {
					errHandlerErrs = append(errHandlerErrs, err)
				}),
			)

			var mu sync.Mutex
			var wg sync.WaitGroup
			wg.Add(len(test.expHashes))
			hashes := make(map[string]bool)
			c.SubscribeHashTx(func(ctx context.Context, hash string) {
				defer wg.Done()
				defer mu.Unlock()
				mu.Lock()
				hashes[hash] = true
			})

			assert.NoError(t, c.Connect())
			wg.Wait()

			assert.Equal(t, len(test.expHashes), len(hashes))
			for _, hash := range test.expHashes {
				assert.True(t, hashes[hash])
			}
		})
	}
}

func TestNodeMQ_SubscribeHashBlock(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		messages  []zmq4.Msg
		expHashes []string
	}{
		"successful messages": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expHashes: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
		},
		"irrelevant messages are ignored": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashTx)
					body, err := hex.DecodeString("b9037b78403cc4e76a06060b4e9d1e1cdf7c85ce7cb6f074e1f9ed2fb6aa10a6")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expHashes: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var errHandlerErrs []error
			socket := &mocks.SocketMock{
				DialFunc: func(addr string) error {
					return nil
				},
				SetOptionFunc: func(opt string, v interface{}) error {
					return nil
				},
				RecvFunc: func() (zmq4.Msg, error) {
					if len(test.messages) == 0 {
						return zmq4.Msg{}, context.Canceled
					}
					defer func() { test.messages = test.messages[1:] }()

					return test.messages[0], nil
				},
				CloseFunc: func() error {
					return nil
				},
			}
			c := zmq.NewNodeMQ(
				zmq.WithHost("tcp://localhost:12345"),
				zmq.WithCustomZMQSocket(socket),
				zmq.WithErrorHandler(func(ctx context.Context, err error) {
					errHandlerErrs = append(errHandlerErrs, err)
				}),
			)

			var mu sync.Mutex
			var wg sync.WaitGroup
			wg.Add(len(test.expHashes))
			hashes := make(map[string]bool, 0)
			assert.NoError(t, c.SubscribeHashBlock(func(ctx context.Context, hash string) {
				defer wg.Done()
				defer mu.Unlock()
				mu.Lock()
				hashes[hash] = true
			}))

			assert.NoError(t, c.Connect())
			wg.Wait()

			assert.Equal(t, len(test.expHashes), len(hashes))
			for _, hash := range test.expHashes {
				assert.True(t, hashes[hash])
			}
		})
	}
}

func TestNodeMQ_SubscribeRawTx(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		messages []zmq4.Msg
		expTxs   []string
		expErrs  []string
	}{
		"successful messages": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawTx)
					body, err := hex.DecodeString("020000000163637a131f20fe7a110db51ff2d3ec4815f344c3c8cc328a5b432c8c286b8f7c00000000484730440220386d4130664137943b157ae9584091ba703f5fc283b6d0f44db035757c017e440220379d94394944d3f01c62501940c1c35b672c755f58e365de1e0c85f14f62d75641feffffff0200e1f505000000001976a914beb20631d5271a6e150231e625bccff55a58cbea88ac40101024010000001976a9148c4a28cfd190444bac5945da342944d6b61e4ae088ac65000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawTx)
					body, err := hex.DecodeString("02000000019da0ccb19dddcf507b0e5b81df8c79f9db7531d09c93b937bec0d1d8e1ff44a4010000006b483045022100e166888a3bfc414111f9d8e4339b4d5e738b995852a7179fd3d8f6eac767fef5022022fcf5f3c8ebe784f7a111dff92adcfd21b42d365249f8335ded08cab6c58cd741210382a6573a2a3253d3264071510045f962aaf4342a996825f831ead59785f5bd0bfeffffff025e2e1a1e010000001976a91465207be9504233b5d9bf57846d9e9a1abdc55b4188ac00e1f505000000001976a914beb20631d5271a6e150231e625bccff55a58cbea88ac66000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expTxs: []string{
				"a444ffe1d8d1c0be37b9939cd03175dbf9798cdf815b0e7b50cfdd9db1cca09d",
				"e272917b10474a4bbf6d760fe0caa443da23d0788989f73bd7c4bc2ff963221d",
			},
		},
		"irrelevant messages are ignored": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawTx)
					body, err := hex.DecodeString("020000000163637a131f20fe7a110db51ff2d3ec4815f344c3c8cc328a5b432c8c286b8f7c00000000484730440220386d4130664137943b157ae9584091ba703f5fc283b6d0f44db035757c017e440220379d94394944d3f01c62501940c1c35b672c755f58e365de1e0c85f14f62d75641feffffff0200e1f505000000001976a914beb20631d5271a6e150231e625bccff55a58cbea88ac40101024010000001976a9148c4a28cfd190444bac5945da342944d6b61e4ae088ac65000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawTx)
					body, err := hex.DecodeString("02000000019da0ccb19dddcf507b0e5b81df8c79f9db7531d09c93b937bec0d1d8e1ff44a4010000006b483045022100e166888a3bfc414111f9d8e4339b4d5e738b995852a7179fd3d8f6eac767fef5022022fcf5f3c8ebe784f7a111dff92adcfd21b42d365249f8335ded08cab6c58cd741210382a6573a2a3253d3264071510045f962aaf4342a996825f831ead59785f5bd0bfeffffff025e2e1a1e010000001976a91465207be9504233b5d9bf57846d9e9a1abdc55b4188ac00e1f505000000001976a914beb20631d5271a6e150231e625bccff55a58cbea88ac66000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expTxs: []string{
				"a444ffe1d8d1c0be37b9939cd03175dbf9798cdf815b0e7b50cfdd9db1cca09d",
				"e272917b10474a4bbf6d760fe0caa443da23d0788989f73bd7c4bc2ff963221d",
			},
		},
		"errors are relayed to the error handler": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawTx)
					body, err := hex.DecodeString("000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawTx)
					body, err := hex.DecodeString("020000000163637a131f20fe7a110db51ff2d3ec4815f344c3c8cc328a5b432c8c286b8f7c00000000484730440220386d4130664137943b157ae9584091ba703f5fc283b6d0f44db035757c017e440220379d94394944d3f01c62501940c1c35b672c755f58e365de1e0c85f14f62d75641feffffff0200e1f505000000001976a914beb20631d5271a6e150231e625bccff55a58cbea88ac40101024010000001976a9148c4a28cfd190444bac5945da342944d6b61e4ae088ac65000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawTx)
					body, err := hex.DecodeString("02000000019da0ccb19dddcf507b0e5b81df8c79f9db7531d09c93b937bec0d1d8e1ff44a4010000006b483045022100e166888a3bfc414111f9d8e4339b4d5e738b995852a7179fd3d8f6eac767fef5022022fcf5f3c8ebe784f7a111dff92adcfd21b42d365249f8335ded08cab6c58cd741210382a6573a2a3253d3264071510045f962aaf4342a996825f831ead59785f5bd0bfeffffff025e2e1a1e010000001976a91465207be9504233b5d9bf57846d9e9a1abdc55b4188ac00e1f505000000001976a914beb20631d5271a6e150231e625bccff55a58cbea88ac66000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expTxs: []string{
				"a444ffe1d8d1c0be37b9939cd03175dbf9798cdf815b0e7b50cfdd9db1cca09d",
				"e272917b10474a4bbf6d760fe0caa443da23d0788989f73bd7c4bc2ff963221d",
			},
			expErrs: []string{
				"nLockTime length must be 4 bytes long",
			},
		},
	}

	for name, test := range tests {
		var msgMu, errMu sync.Mutex
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(len(test.expTxs) + len(test.expErrs))

			socket := &mocks.SocketMock{
				DialFunc: func(addr string) error {
					return nil
				},
				SetOptionFunc: func(opt string, v interface{}) error {
					return nil
				},
				RecvFunc: func() (zmq4.Msg, error) {
					if len(test.messages) == 0 {
						return zmq4.Msg{}, context.Canceled
					}
					defer func() { test.messages = test.messages[1:] }()

					return test.messages[0], nil
				},
				CloseFunc: func() error {
					return nil
				},
			}
			var errs []string
			c := zmq.NewNodeMQ(
				zmq.WithHost("tcp://localhost:12345"),
				zmq.WithRaw(),
				zmq.WithCustomZMQSocket(socket),
				zmq.WithErrorHandler(func(ctx context.Context, err error) {
					defer wg.Done()
					defer errMu.Unlock()
					errMu.Lock()
					errs = append(errs, err.Error())
				}),
			)
			txs := make(map[string]bool, 0)
			assert.NoError(t, c.SubscribeRawTx(func(ctx context.Context, tx *bt.Tx) {
				defer wg.Done()
				defer msgMu.Unlock()
				msgMu.Lock()
				txs[tx.TxID()] = true
			}))

			assert.NoError(t, c.Connect())
			wg.Wait()

			assert.Equal(t, len(test.expTxs), len(txs))
			for _, tx := range test.expTxs {
				assert.True(t, txs[tx])
			}
			assert.Equal(t, test.expErrs, errs)
		})
	}
}

func TestNodeMQ_SubscribeRawBlock(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		messages  []zmq4.Msg
		expBlocks []string
		expErrs   []string
	}{
		"successful messages": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawBlock)
					body, err := hex.DecodeString("000000202e786df0a69bb12816d02d7017bd7ab1137ba1674f1c7c35cd216a96bcb3df1f59ae405b34f3beb08e1cf6475a9e0a19a63b29127340172ab7ceeef7c4031df1042d9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401680101ffffffff0100f2052a01000000232103d427b882f2ec8d379353e478061a962135830501ed5d06a58af03ad5c311dfbdac00000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawBlock)
					body, err := hex.DecodeString("000000201690a0f7650650aa3ba965f9703b94dcb471307dcffe2c2cc116f4f8e5360b54a4066d10b4bf88e9f4d95df3c788bc024e005fb79cce357f1701b8a7a4fdbe124d2e9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401690101ffffffff0100f2052a01000000232103c3a2396bdfa58e93f8bb80f252aabf591d827a701988119f3e6c947c230d25d3ac00000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expBlocks: []string{
				"000000202e786df0a69bb12816d02d7017bd7ab1137ba1674f1c7c35cd216a96bcb3df1f59ae405b34f3beb08e1cf6475a9e0a19a63b29127340172ab7ceeef7c4031df1042d9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401680101ffffffff0100f2052a01000000232103d427b882f2ec8d379353e478061a962135830501ed5d06a58af03ad5c311dfbdac00000000",
				"000000201690a0f7650650aa3ba965f9703b94dcb471307dcffe2c2cc116f4f8e5360b54a4066d10b4bf88e9f4d95df3c788bc024e005fb79cce357f1701b8a7a4fdbe124d2e9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401690101ffffffff0100f2052a01000000232103c3a2396bdfa58e93f8bb80f252aabf591d827a701988119f3e6c947c230d25d3ac00000000",
			},
		},
		"irrelevant messages are ignored": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashTx)
					body, err := hex.DecodeString("a444ffe1d8d1c0be37b9939cd03175dbf9798cdf815b0e7b50cfdd9db1cca09d")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawBlock)
					body, err := hex.DecodeString("000000202e786df0a69bb12816d02d7017bd7ab1137ba1674f1c7c35cd216a96bcb3df1f59ae405b34f3beb08e1cf6475a9e0a19a63b29127340172ab7ceeef7c4031df1042d9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401680101ffffffff0100f2052a01000000232103d427b882f2ec8d379353e478061a962135830501ed5d06a58af03ad5c311dfbdac00000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawBlock)
					body, err := hex.DecodeString("000000201690a0f7650650aa3ba965f9703b94dcb471307dcffe2c2cc116f4f8e5360b54a4066d10b4bf88e9f4d95df3c788bc024e005fb79cce357f1701b8a7a4fdbe124d2e9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401690101ffffffff0100f2052a01000000232103c3a2396bdfa58e93f8bb80f252aabf591d827a701988119f3e6c947c230d25d3ac00000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expBlocks: []string{
				"000000202e786df0a69bb12816d02d7017bd7ab1137ba1674f1c7c35cd216a96bcb3df1f59ae405b34f3beb08e1cf6475a9e0a19a63b29127340172ab7ceeef7c4031df1042d9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401680101ffffffff0100f2052a01000000232103d427b882f2ec8d379353e478061a962135830501ed5d06a58af03ad5c311dfbdac00000000",
				"000000201690a0f7650650aa3ba965f9703b94dcb471307dcffe2c2cc116f4f8e5360b54a4066d10b4bf88e9f4d95df3c788bc024e005fb79cce357f1701b8a7a4fdbe124d2e9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401690101ffffffff0100f2052a01000000232103c3a2396bdfa58e93f8bb80f252aabf591d827a701988119f3e6c947c230d25d3ac00000000",
			},
		},
		"errors are relayed to the error handler": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawBlock)
					body, err := hex.DecodeString("000000202e786df0a69bb12816d02d7017bd7ab1137ba1674f1c7c35cd216a96bcb3df1f59ae405b34f3beb08e1cf6475a9e0a19a63b29127340172ab7ceeef7c4031df1042d9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401680101ffffffff0100f2052a01000000232103d427b882f2ec8d379353e478061a962135830500")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawBlock)
					body, err := hex.DecodeString("000000202e786df0a69bb12816d02d7017bd7ab1137ba1674f1c7c35cd216a96bcb3df1f59ae405b34f3beb08e1cf6475a9e0a19a63b29127340172ab7ceeef7c4031df1042d9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401680101ffffffff0100f2052a01000000232103d427b882f2ec8d379353e478061a962135830501ed5d06a58af03ad5c311dfbdac00000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRawBlock)
					body, err := hex.DecodeString("000000201690a0f7650650aa3ba965f9703b94dcb471307dcffe2c2cc116f4f8e5360b54a4066d10b4bf88e9f4d95df3c788bc024e005fb79cce357f1701b8a7a4fdbe124d2e9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401690101ffffffff0100f2052a01000000232103c3a2396bdfa58e93f8bb80f252aabf591d827a701988119f3e6c947c230d25d3ac00000000")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}},
			expBlocks: []string{
				"000000202e786df0a69bb12816d02d7017bd7ab1137ba1674f1c7c35cd216a96bcb3df1f59ae405b34f3beb08e1cf6475a9e0a19a63b29127340172ab7ceeef7c4031df1042d9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401680101ffffffff0100f2052a01000000232103d427b882f2ec8d379353e478061a962135830501ed5d06a58af03ad5c311dfbdac00000000",
				"000000201690a0f7650650aa3ba965f9703b94dcb471307dcffe2c2cc116f4f8e5360b54a4066d10b4bf88e9f4d95df3c788bc024e005fb79cce357f1701b8a7a4fdbe124d2e9261ffff7f20000000000102000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401690101ffffffff0100f2052a01000000232103c3a2396bdfa58e93f8bb80f252aabf591d827a701988119f3e6c947c230d25d3ac00000000",
			},
			expErrs: []string{
				"input length too short < 8 + script",
			},
		},
	}

	for name, test := range tests {
		var msgMu, errMu sync.Mutex
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(len(test.expBlocks) + len(test.expErrs))

			socket := &mocks.SocketMock{
				DialFunc: func(addr string) error {
					return nil
				},
				SetOptionFunc: func(opt string, v interface{}) error {
					return nil
				},
				RecvFunc: func() (zmq4.Msg, error) {
					if len(test.messages) == 0 {
						return zmq4.Msg{}, context.Canceled
					}
					defer func() { test.messages = test.messages[1:] }()

					return test.messages[0], nil
				},
				CloseFunc: func() error {
					return nil
				},
			}
			var errs []string
			c := zmq.NewNodeMQ(
				zmq.WithHost("tcp://localhost:12345"),
				zmq.WithRaw(),
				zmq.WithCustomZMQSocket(socket),
				zmq.WithErrorHandler(func(ctx context.Context, err error) {
					defer wg.Done()
					defer errMu.Unlock()
					errMu.Lock()
					errs = append(errs, err.Error())
				}),
			)
			blocks := make(map[string]bool, 0)
			assert.NoError(t, c.SubscribeRawBlock(func(ctx context.Context, blk *bc.Block) {
				defer wg.Done()
				defer msgMu.Unlock()
				msgMu.Lock()
				blocks[blk.String()] = true
			}))

			assert.NoError(t, c.Connect())
			wg.Wait()

			assert.Equal(t, len(test.expBlocks), len(blocks))
			for _, blk := range test.expBlocks {
				assert.True(t, blocks[blk])
			}
			assert.Equal(t, test.expErrs, errs)
		})
	}
}

func TestNodeMQ_SubscribeDiscardFromMempool(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		messages    []zmq4.Msg
		expDiscards []string
		expErrs     []string
	}{
		"successful messages": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicDiscardFromMempool)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicDiscardFromMempool)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64"}`)

					return [][]byte{header, body}
				}(),
			}},
			expDiscards: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
		},
		"irrelevant messages are ignored": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicDiscardFromMempool)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicDiscardFromMempool)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64"}`)

					return [][]byte{header, body}
				}(),
			}},
			expDiscards: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
		},
		"errors are relayed to the error handler": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicDiscardFromMempool)
					body := []byte(`{"txid":"abc123","bloc"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicDiscardFromMempool)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicDiscardFromMempool)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64"}`)

					return [][]byte{header, body}
				}(),
			}},
			expDiscards: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
			expErrs: []string{
				"invalid character '}' after object key",
			},
		},
	}

	for name, test := range tests {
		var msgMu, errMu sync.Mutex
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(len(test.expDiscards) + len(test.expErrs))

			socket := &mocks.SocketMock{
				DialFunc: func(addr string) error {
					return nil
				},
				SetOptionFunc: func(opt string, v interface{}) error {
					return nil
				},
				RecvFunc: func() (zmq4.Msg, error) {
					if len(test.messages) == 0 {
						return zmq4.Msg{}, context.Canceled
					}
					defer func() { test.messages = test.messages[1:] }()

					return test.messages[0], nil
				},
				CloseFunc: func() error {
					return nil
				},
			}
			var errs []string
			c := zmq.NewNodeMQ(
				zmq.WithHost("tcp://localhost:12345"),
				zmq.WithRaw(),
				zmq.WithCustomZMQSocket(socket),
				zmq.WithErrorHandler(func(ctx context.Context, err error) {
					defer wg.Done()
					defer errMu.Unlock()
					errMu.Lock()
					errs = append(errs, err.Error())
				}),
			)
			discards := make(map[string]bool, 0)
			assert.NoError(t, c.SubscribeDiscardFromMempool(func(ctx context.Context, msg *zmq.MempoolDiscard) {
				defer wg.Done()
				defer msgMu.Unlock()
				msgMu.Lock()
				discards[msg.BlockHash] = true
			}))

			assert.NoError(t, c.Connect())
			wg.Wait()

			assert.Equal(t, len(test.expDiscards), len(discards))
			for _, discard := range test.expDiscards {
				assert.True(t, discards[discard])
			}
			assert.Equal(t, test.expErrs, errs)
		})
	}
}

func TestNodeMQ_SubscribeRemovedFromMempoolBlock(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		messages    []zmq4.Msg
		expRemovals []string
		expErrs     []string
	}{
		"successful messages": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRemovedFromMempoolBlock)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRemovedFromMempoolBlock)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64"}`)

					return [][]byte{header, body}
				}(),
			}},
			expRemovals: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
		},
		"irrelevant messages are ignored": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicHashBlock)
					body, err := hex.DecodeString("000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af")
					assert.NoError(t, err)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRemovedFromMempoolBlock)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRemovedFromMempoolBlock)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64"}`)

					return [][]byte{header, body}
				}(),
			}},
			expRemovals: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
		},
		"errors are relayed to the error handler": {
			messages: []zmq4.Msg{{
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRemovedFromMempoolBlock)
					body := []byte(`{"txid":"abc123","bloc"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRemovedFromMempoolBlock)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af"}`)

					return [][]byte{header, body}
				}(),
			}, {
				Frames: func() [][]byte {
					header := []byte(zmq.TopicRemovedFromMempoolBlock)
					body := []byte(`{"txid":"abc123","blockhash":"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64"}`)

					return [][]byte{header, body}
				}(),
			}},
			expRemovals: []string{
				"000000000000000001cd535a5b3ad0fb3ec22d153e845508666818ab29eb27af",
				"000000000000000002d8a4faad6d8c026526bf1a8a2abe2074a42e71f95ebb64",
			},
			expErrs: []string{
				"invalid character '}' after object key",
			},
		},
	}

	for name, test := range tests {
		var msgMu, errMu sync.Mutex
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(len(test.expRemovals) + len(test.expErrs))

			socket := &mocks.SocketMock{
				DialFunc: func(addr string) error {
					return nil
				},
				SetOptionFunc: func(opt string, v interface{}) error {
					return nil
				},
				RecvFunc: func() (zmq4.Msg, error) {
					if len(test.messages) == 0 {
						return zmq4.Msg{}, context.Canceled
					}
					defer func() { test.messages = test.messages[1:] }()

					return test.messages[0], nil
				},
				CloseFunc: func() error {
					return nil
				},
			}
			var errs []string
			c := zmq.NewNodeMQ(
				zmq.WithHost("tcp://localhost:12345"),
				zmq.WithRaw(),
				zmq.WithCustomZMQSocket(socket),
				zmq.WithErrorHandler(func(ctx context.Context, err error) {
					defer wg.Done()
					defer errMu.Unlock()
					errMu.Lock()
					errs = append(errs, err.Error())
				}),
			)
			discards := make(map[string]bool, 0)
			assert.NoError(t, c.SubscribeRemovedFromMempoolBlock(func(ctx context.Context, msg *zmq.MempoolDiscard) {
				defer wg.Done()
				defer msgMu.Unlock()
				msgMu.Lock()
				discards[msg.BlockHash] = true
			}))

			assert.NoError(t, c.Connect())
			wg.Wait()

			assert.Equal(t, len(test.expRemovals), len(discards))
			for _, discard := range test.expRemovals {
				assert.True(t, discards[discard])
			}
			assert.Equal(t, test.expErrs, errs)
		})
	}
}
