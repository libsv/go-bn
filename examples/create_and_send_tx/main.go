package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libsv/go-bn"
	"github.com/libsv/go-bt/v2"
)

func main() {
	c := bn.NewNodeClient(
		bn.WithHost("http://localhost:18332"),
		bn.WithCreds("bitcoin", "bitcoin"),
	)
	ctx := context.Background()

	tx := bt.NewTx()
	if err := tx.AddP2PKHOutputFromAddress("n4LvK5SVxp8ohxLHS6fXz47ErBLg5WDHgS", 20000); err != nil {
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

	txID, err := c.SendRawTransaction(ctx, signedTx.Tx, nil)
	if err != nil {
		panic(err)
	}

	resp, err := c.RawTransaction(ctx, txID)
	if err != nil {
		panic(err)
	}

	bb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bb))
}
