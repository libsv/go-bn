package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libsv/go-bn"
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

			blockHeader, err := c.BlockHeader(ctx, "2992bf02b5ee83676bcd6e19aac395ce624cfed74ad1a2d9875789bdb3aab194")
			if err != nil {
				panic(err)
			}

			bb, err := json.MarshalIndent(blockHeader, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(bb))
		}()
	}

	wg.Wait()
}
