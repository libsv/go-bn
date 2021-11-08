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

			resp, err := c.WalletInfo(ctx)
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
