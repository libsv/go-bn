package bn_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/libsv/go-bn"
	"github.com/libsv/go-bn/internal/config"
	"github.com/libsv/go-bn/internal/mocks"
	"github.com/libsv/go-bn/internal/service"
	"github.com/libsv/go-bn/models"
	"github.com/libsv/go-bn/testing/util"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/sighash"
	"github.com/stretchr/testify/assert"
)

func TestTxClient_CreateRawTransaction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testFile   string
		utxos      bt.UTXOs
		params     models.ParamsCreateRawTransaction
		expParams  map[string]float64
		expRequest models.Request
		expTx      string
		expErr     error
	}{
		"successful query": {
			testFile: "createrawtx",
			utxos:    bt.UTXOs{},
			expTx:    "02000000000100e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000",
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "createrawtransaction",
				Params:  []interface{}{[]interface{}{}, map[string]interface{}{"mpzLdVLZhbRXxYpaT8YcHntWb2tyPJvUnz": 0.1}},
			},
			params: models.ParamsCreateRawTransaction{
				Outputs: func() []*bt.Output {
					tx := bt.NewTx()
					assert.NoError(t, tx.AddP2PKHOutputFromAddress("mpzLdVLZhbRXxYpaT8YcHntWb2tyPJvUnz", 10000000))
					return tx.Outputs
				}(),
			},
			expParams: map[string]float64{
				"mpzLdVLZhbRXxYpaT8YcHntWb2tyPJvUnz": 0.1,
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

			c := bn.NewTransactionClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "createrawtransaction", method)
						assert.Equal(t, 2, len(args))
						assert.Equal(t, args[0], test.utxos.NodeJSON())
						assert.Equal(t, args[1], test.expParams)

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			tx, err := c.CreateRawTransaction(context.TODO(), test.utxos, test.params)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expTx, tx.String())
			}
		})
	}
}

func TestTxClient_FundRawTransaction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testFile    string
		tx          *bt.Tx
		opts        *models.OptsFundRawTransaction
		expRequest  models.Request
		expResponse *models.FundRawTransaction
		expTx       string
		expErr      error
	}{
		"successful query with opts": {
			testFile: "fundrawtx",
			tx: func() *bt.Tx {
				tx, err := bt.NewTxFromString("02000000000100e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
				assert.NoError(t, err)
				return tx
			}(),
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "fundrawtransaction",
				Params: []interface{}{
					"02000000000100e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000",
					map[string]interface{}{
						"changeAddress":    "wow",
						"changePosition":   1.0,
						"feeRate":          5.0,
						"includeWatching":  true,
						"reserveChangeKey": true,
					}},
			},
			expResponse: &models.FundRawTransaction{
				Fee:            226,
				ChangePosition: 0,
				Tx: func() *bt.Tx {
					tx, err := bt.NewTxFromString("0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc90000000000feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
					assert.NoError(t, err)
					return tx
				}(),
			},
			opts: &models.OptsFundRawTransaction{
				ChangeAddress:    "wow",
				ChangePosition:   1,
				FeeRate:          5,
				IncludeWatching:  true,
				LockUnspents:     false,
				ReserveChangeKey: func() *bool { s := true; return &s }(),
			},
		},
		"successful query no opts": {
			testFile: "fundrawtx",
			tx: func() *bt.Tx {
				tx, err := bt.NewTxFromString("02000000000100e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
				assert.NoError(t, err)
				return tx
			}(),
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "fundrawtransaction",
				Params: []interface{}{
					"02000000000100e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000",
				},
			},
			expResponse: &models.FundRawTransaction{
				Fee:            226,
				ChangePosition: 0,
				Tx: func() *bt.Tx {
					tx, err := bt.NewTxFromString("0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc90000000000feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
					assert.NoError(t, err)
					return tx
				}(),
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

			c := bn.NewTransactionClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "fundrawtransaction", method)
						if test.opts == nil {
							assert.Equal(t, 1, len(args))
						} else {
							assert.Equal(t, 2, len(args))
							assert.Equal(t, test.opts, args[1])
						}

						assert.Equal(t, test.tx.String(), args[0])

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			fundedTx, err := c.FundRawTransaction(context.TODO(), test.tx, test.opts)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expResponse, fundedTx)
			}
		})
	}
}

func TestTxClient_RawTransaction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testFile   string
		txID       string
		expRequest models.Request
		expTx      string
		expErr     error
	}{
		"successful query": {
			testFile: "getrawtx",
			txID:     "c98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fb",
			expTx:    "0200000001c9059cca32a90834a9ea6e989446edb4282e91bba486f4512477052214b185df0000000048473044022056e7348677c69dbcba776fbe0c270116c2a3eaf0bead0c1ccdbd9c083b73a08e022062da00341e54a28bb83b28dfd772c9504f5aace3452e762dc30dff249a378c0a41feffffff0240101024010000001976a914316230517501a16e2837465ec28c157fa61cabec88ac00e1f505000000001976a914beb20631d5271a6e150231e625bccff55a58cbea88ac70000000",
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "getrawtransaction",
				Params:  []interface{}{"c98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fb", true},
			},
		},
		"error is reported": {
			txID:     "c98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fc",
			testFile: "getrawtx_notfound",
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "getrawtransaction",
				Params:  []interface{}{"c98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fc", true},
			},
			expErr: errors.New("-5: No such mempool or blockchain transaction. Use gettransaction for wallet transactions."),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewTransactionClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "getrawtransaction", method)
						assert.Equal(t, 2, len(args))
						assert.Equal(t, args[0], test.txID)
						assert.Equal(t, args[1], true)

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			tx, err := c.RawTransaction(context.TODO(), test.txID)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expTx, tx.String())
			}
		})
	}
}

func TestTxClient_SignRawTransaction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testFile    string
		tx          *bt.Tx
		opts        *models.OptsSignRawTransaction
		expRequest  models.Request
		expResponse *models.SignedRawTransaction
		expTx       string
		expErr      error
	}{
		"successful query with opts": {
			testFile: "signrawtx",
			tx: func() *bt.Tx {
				tx, err := bt.NewTxFromString("0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc90000000000feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
				assert.NoError(t, err)
				return tx
			}(),
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "signrawtransaction",
				Params: []interface{}{
					"0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc90000000000feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000",
					[]interface{}{}, []interface{}{"myprivkey"}, "ALL|FORKID|ANYONECANPAY",
				},
			},
			expResponse: &models.SignedRawTransaction{
				Tx: func() *bt.Tx {
					tx, err := bt.NewTxFromString("0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc9000000006b483045022100d50174438859f148a9f21dfc98a7e3d51a010f279513a3ecb6375d2f10e4676102201668d8ca301d8d0cc28d077ce5661cb815b9f7df518ef3c741f815639cf5ba784121034df56fcde16931d7059669da5fa8ae845aab89bc7b3f9e6cbe2b3f7322315389feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
					assert.NoError(t, err)
					return tx
				}(),
				Complete: true,
			},
			opts: &models.OptsSignRawTransaction{
				From:        bt.UTXOs{},
				PrivateKeys: []string{"myprivkey"},
				SigHashType: sighash.AllForkID | sighash.AnyOneCanPayForkID,
			},
		},
		"successful query no opts": {
			testFile: "signrawtx",
			tx: func() *bt.Tx {
				tx, err := bt.NewTxFromString("0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc90000000000feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
				assert.NoError(t, err)
				return tx
			}(),
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "signrawtransaction",
				Params: []interface{}{
					"0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc90000000000feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000",
				},
			},
			expResponse: &models.SignedRawTransaction{
				Tx: func() *bt.Tx {
					tx, err := bt.NewTxFromString("0200000001fbb877c83aaf682f74611628b0088254c8094fff9cf6328ed969c587112b8fc9000000006b483045022100d50174438859f148a9f21dfc98a7e3d51a010f279513a3ecb6375d2f10e4676102201668d8ca301d8d0cc28d077ce5661cb815b9f7df518ef3c741f815639cf5ba784121034df56fcde16931d7059669da5fa8ae845aab89bc7b3f9e6cbe2b3f7322315389feffffff025e2e1a1e010000001976a91401becd83278806a62cd87bed129faa72af38a0d588ac00e1f505000000001976a91467e701e630adaee761583a894b53d4356028ca0b88ac00000000")
					assert.NoError(t, err)
					return tx
				}(),
				Complete: true,
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

			c := bn.NewTransactionClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "signrawtransaction", method)
						assert.Equal(t, test.tx.String(), args[0])
						if test.opts == nil {
							assert.Equal(t, 1, len(args))
						} else {
							assert.Equal(t, 4, len(args))
							assert.Equal(t, test.opts.From.NodeJSON(), args[1])
							assert.Equal(t, test.opts.PrivateKeys, args[2])
							assert.Equal(t, test.opts.SigHashType.String(), args[3])
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			signedTx, err := c.SignRawTransaction(context.TODO(), test.tx, test.opts)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expResponse, signedTx)
			}
		})
	}
}
