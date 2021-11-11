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

func TestUtilClient_SignMessageWithPrivKey(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		wif        *wif.WIF
		msg        string
		expRequest models.Request
		expMsg     string
		expErr     error
	}{
		"successful request": {
			testFile: "signmessagewithprivkey",
			wif: func() *wif.WIF {
				wif, err := wif.DecodeWIF("cW9n4pgq9MqqGD8Ux5cwpgJAJ1VzPvZgskbCEmK1QmWUicejRFQn")
				assert.NoError(t, err)
				return wif
			}(),
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "signmessagewithprivkey",
				Params:  []interface{}{"cW9n4pgq9MqqGD8Ux5cwpgJAJ1VzPvZgskbCEmK1QmWUicejRFQn", "hello"},
			},
			msg:    "hello",
			expMsg: "IL4oekQr7n8+u6QWCvZ+jMFhRz/zMMq4wfBvXhh+eP/zVzknU+IteOsEwyGguMnN/m7BvtOdf5b9JofdI4jEktI=",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewUtilClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "signmessagewithprivkey", method)
						assert.Equal(t, len(args), 2)
						assert.Equal(t, test.wif.String(), args[0])
						assert.Equal(t, test.msg, args[1])

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			signedMsg, err := c.SignMessageWithPrivKey(context.TODO(), test.wif, test.msg)
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expMsg, signedMsg)
			}
		})
	}
}
