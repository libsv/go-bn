package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libsv/go-bn"
	"github.com/libsv/go-bn/internal/util"
	"github.com/libsv/go-bt/v2"
)

func main() {
	c := bn.NewNodeClient(
		bn.WithHost("http://localhost:18332"),
		bn.WithCreds("bitcoin", "bitcoin"),
	)
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			tx := bt.NewTx()
			if err := tx.AddP2PKHOutputFromAddress("n4LvK5SVxp8ohxLHS6fXz47ErBLg5WDHgS", util.SatoshisToBSV(2)); err != nil {
				panic(err)
			}
			fundedTx, err := c.FundRawTransaction(ctx, tx, nil)
			if err != nil {
				panic(err)
			}
			signedTx, err := c.SignRawTransaction(ctx, fundedTx.Tx, nil)
			if err != nil {
				panic(err)
			}

			bb, err := json.MarshalIndent(signedTx.Tx, "", "  ")
			fmt.Println(string(bb))
		}()
	}

	wg.Wait()
}
