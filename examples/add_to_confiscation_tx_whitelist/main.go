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

	funds := []models.Fund{
		{
			TxOut: models.TxOut{
				TxId: "",
				Vout: 0,
			},
			EnforceAtHeight: []models.Enforce{
				{
					Start: 100000,
					Stop:  100001,
				},
			},
			PolicyExpiresWithConsensus: false,
		},
	}

	resp, err := c.AddToConsensusBlacklist(ctx, funds)
	if err != nil {
		panic(err)
	}

	bb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bb))
}
