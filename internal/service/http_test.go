package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/libsv/go-bn/internal/config"
	"github.com/libsv/go-bn/internal/service"
	"github.com/libsv/go-bn/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestRPC_Do_SingleFlight(t *testing.T) {
	t.Parallel()
	type invocation struct {
		timesCalled int32
		method      string
		args        []interface{}
	}
	tests := map[string]struct {
		invocations []invocation
		expCalls    int32
		expErr      error
	}{
		"single flight same data": {
			expCalls: 1,
			invocations: []invocation{{
				timesCalled: 500,
				method:      "getrawtransaction",
				args:        []interface{}{"c98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fb"},
			}},
		},
		"single flight diff args same method": {
			expCalls: 4,
			invocations: []invocation{{
				timesCalled: 500,
				method:      "getrawtransaction",
				args:        []interface{}{"c98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fb"},
			}, {
				timesCalled: 500,
				method:      "getrawtransaction",
				args:        []interface{}{"t98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fz"},
			}, {
				timesCalled: 400,
				method:      "getrawtransaction",
				args:        []interface{}{"s98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fx"},
			}, {
				timesCalled: 600,
				method:      "getrawtransaction",
				args:        []interface{}{"a98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fa"},
			}},
		},
		"single flight diff method same args": {
			expCalls: 3,
			invocations: []invocation{{
				timesCalled: 400,
				method:      "getinfo",
				args:        []interface{}{},
			}, {
				timesCalled: 500,
				method:      "getrawtransaction",
				args:        []interface{}{"a98f2b1187c569d98e32f69cff4f09c8548208b0281661742f68af3ac877b8fa"},
			}, {
				timesCalled: 1000,
				method:      "getmininginfo",
				args:        []interface{}{},
			}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var timesCalled int32
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&timesCalled, 1)
				time.Sleep(1 * time.Second)

				bb, err := json.Marshal(models.Response{
					Result: "ohiya",
				})
				assert.NoError(t, err)
				w.Write(bb)
			}))
			defer svr.Close()

			c := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			g, ctx := errgroup.WithContext(context.TODO())
			for _, inv := range test.invocations {
				inv := inv
				g.Go(func() error {
					for i := 0; i < int(inv.timesCalled); i++ {
						g.Go(func() error { return c.Do(ctx, inv.method, nil, inv.args...) })
					}
					return nil
				})
			}

			err := g.Wait()
			if test.expErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expCalls, timesCalled)
		})
	}
}
