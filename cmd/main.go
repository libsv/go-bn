package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libsv/go-bn"
	"github.com/libsv/go-bn/models"
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

			utxos, err := c.ListUnspent(ctx, nil)
			if err != nil {
				panic(err)
			}

			tx, err := c.CreateRawTransaction(ctx, utxos[:1], models.ParamsCreateRawTransaction{
				Outputs: func() []*bt.Output {
					tx := bt.NewTx()
					if err = tx.AddP2PKHOutputFromAddress("n4LvK5SVxp8ohxLHS6fXz47ErBLg5WDHgS", utxos[0].Satoshis/2); err != nil {
						panic(err)
					}
					return tx.Outputs
				}(),
			})
			if err != nil {
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

			tx2, err := c.CreateRawTransaction(ctx, utxos[1:2], models.ParamsCreateRawTransaction{
				Outputs: func() []*bt.Output {
					tx := bt.NewTx()
					if err = tx.AddP2PKHOutputFromAddress("n4LvK5SVxp8ohxLHS6fXz47ErBLg5WDHgS", utxos[0].Satoshis/2); err != nil {
						panic(err)
					}
					return tx.Outputs
				}(),
			})
			if err != nil {
				panic(err)
			}
			fundedTx2, err := c.FundRawTransaction(ctx, tx2, nil)
			if err != nil {
				panic(err)
			}
			signedTx2, err := c.SignRawTransaction(ctx, fundedTx2.Tx, nil)
			if err != nil {
				panic(err)
			}

			resp, err := c.SendRawTransactions(ctx, models.ParamsSendRawTransactions{
				Hex: signedTx.Tx.String(),
			}, models.ParamsSendRawTransactions{
				Hex: signedTx2.Tx.String(),
			})
			if err != nil {
				panic(err)
			}

			bb, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(bb))
		}()
	}

	wg.Wait()
}
