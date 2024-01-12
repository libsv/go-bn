package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libsv/go-bn/models"

	"github.com/libsv/go-bn"
)

func main() {
	c := bn.NewNodeClient(
		bn.WithHost("http://localhost:8333"),
		bn.WithCreds("galt", "galt"),
	)
	ctx := context.Background()

	txs := []models.ConfiscationTransactionDetails{
		{
			ConfiscationTransaction: models.ConfiscationTransaction{
				EnforceAtHeight: 10000,
				Hex:             "",
			},
		},
	}

	resp, err := c.AddToConfiscationTransactionWhitelist(ctx, txs)
	if err != nil {
		panic(err)
	}

	bb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bb))
}
