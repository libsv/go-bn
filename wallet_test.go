package bn_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bn"
	"github.com/libsv/go-bn/internal/config"
	"github.com/libsv/go-bn/internal/mocks"
	"github.com/libsv/go-bn/internal/service"
	"github.com/libsv/go-bn/models"
	"github.com/libsv/go-bn/testing/util"
	"github.com/stretchr/testify/assert"
)

func TestWalletClient_DumpPrivateKey(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		address    string
		expWif     *wif.WIF
		expRequest models.Request
		expErr     error
	}{
		"success request": {
			testFile: "dumpprivkey",
			address:  "mzcEDt2d7QwHazAwD11WWSn8eSCb4gtpSY",
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "dumpprivkey",
				Params:  []interface{}{"mzcEDt2d7QwHazAwD11WWSn8eSCb4gtpSY"},
			},
			expWif: func() *wif.WIF {
				wif, err := wif.DecodeWIF("cW9n4pgq9MqqGD8Ux5cwpgJAJ1VzPvZgskbCEmK1QmWUicejRFQn")
				assert.NoError(t, err)
				return wif
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "dumpprivkey", method)
						assert.Equal(t, 1, len(args))
						assert.Equal(t, test.address, args[0])

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			wif, err := c.DumpPrivateKey(context.TODO(), test.address)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expWif, wif)
			}
		})
	}
}

func TestWalletClient_NewAddress(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		opts       *models.OptsNewAddress
		expRequest models.Request
		expAddress string
		expArgsLen int
		expErr     error
	}{
		"successful request without opts": {
			testFile:   "getnewaddress",
			expAddress: "mxokrvSv54CTNer4Am8WTutjqJcpGS3Txz",
			expArgsLen: 0,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "getnewaddress",
			},
		},
		"successful request with opts": {
			testFile:   "getnewaddress",
			expAddress: "mxokrvSv54CTNer4Am8WTutjqJcpGS3Txz",
			expArgsLen: 1,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "getnewaddress",
				Params:  []interface{}{"accountname"},
			},
			opts: &models.OptsNewAddress{
				Account: "accountname",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "getnewaddress", method)
						assert.Equal(t, test.expArgsLen, len(args))
						if test.opts != nil {
							assert.Equal(t, test.opts.Account, args[0])
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			address, err := c.NewAddress(context.TODO(), test.opts)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expAddress, address)
			}
		})
	}
}

func TestWalletClient_ListAccounts(t *testing.T) {
	tests := map[string]struct {
		testFile    string
		opts        *models.OptsListAccounts
		expRequest  models.Request
		expAccounts map[string]uint64
		expArgsLen  int
		expErr      error
	}{
		"successful request without opts": {
			testFile:   "listaccounts",
			expArgsLen: 0,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "listaccounts",
			},
			expAccounts: map[string]uint64{
				"":     8567000000,
				"john": 100000,
				"bob":  100000000,
			},
		},
		"successful request with opts": {
			testFile:   "listaccounts",
			expArgsLen: 2,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "listaccounts",
				Params:  []interface{}{4.0, true},
			},
			expAccounts: map[string]uint64{
				"":     8567000000,
				"john": 100000,
				"bob":  100000000,
			},
			opts: &models.OptsListAccounts{
				MinConf:          4,
				IncludeWatchOnly: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "listaccounts", method)
						assert.Equal(t, test.expArgsLen, len(args))
						if test.opts != nil {
							assert.Equal(t, test.opts.MinConf, args[0])
							assert.Equal(t, test.opts.IncludeWatchOnly, args[1])
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			accounts, err := c.ListAccounts(context.TODO(), test.opts)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expAccounts, accounts)
			}
		})
	}
}

func TestWalletClient_Move(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		opts       *models.OptsMove
		expRequest models.Request
		expResult  bool
		expArgsLen int
		from       string
		to         string
		amount     uint64
		expAmount  float64
		expErr     error
	}{
		"successful request without opts": {
			testFile:   "move",
			amount:     123456789994,
			from:       "john",
			to:         "bob",
			expResult:  true,
			expArgsLen: 3,
			expAmount:  1234.56789994,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "move",
				Params:  []interface{}{"john", "bob", 1234.56789994},
			},
		},
		"successful request with opts": {
			testFile:   "move",
			amount:     123456789994,
			from:       "john",
			to:         "bob",
			expResult:  true,
			expArgsLen: 5,
			expAmount:  1234.56789994,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "move",
				Params:  []interface{}{"john", "bob", 1234.56789994, "", "oh wow"},
			},
			opts: &models.OptsMove{
				Comment: "oh wow",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "move", method)
						assert.Equal(t, test.expArgsLen, len(args))
						assert.Equal(t, test.from, args[0])
						assert.Equal(t, test.to, args[1])
						assert.Equal(t, test.expAmount, args[2])
						if test.opts != nil {
							assert.Equal(t, test.opts.Comment, args[4])
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			result, err := c.Move(context.TODO(), test.from, test.to, test.amount, test.opts)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expResult, result)
			}
		})
	}
}
