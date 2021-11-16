package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-bn/zmq"
	"github.com/libsv/go-bt/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	z := zmq.NewNodeMQ(
		zmq.WithContext(ctx),
		zmq.WithHost("tcp://localhost:28332"),
		zmq.WithRaw(),
		zmq.WithErrorHandler(func(_ context.Context, err error) {
			fmt.Println("OH NO", err)
		}),
	)

	if err := z.Subscribe(zmq.TopicInvalidTx, func(_ context.Context, bb [][]byte) {
		fmt.Println("invalid tx", hex.EncodeToString(bb[1]))
	}); err != nil {
		panic(err)
	}

	if err := z.SubscribeRawTx(func(_ context.Context, tx *bt.Tx) {
		bb, err := json.Marshal(tx)
		if err != nil {
			panic(err)
		}

		fmt.Println("TX:", string(bb))
	}); err != nil {
		panic(err)
	}

	if err := z.SubscribeRawBlock(func(_ context.Context, blk *bc.Block) {
		bb, err := json.Marshal(blk)
		if err != nil {
			panic(err)
		}

		fmt.Println("Block:", string(bb))
	}); err != nil {
		panic(err)
	}

	if err := z.SubscribeHashBlock(func(_ context.Context, hash string) {
		fmt.Println("BLOCK HASH", hash)
	}); err != nil {
		panic(err)
	}

	log.Fatal(z.Connect())
}
