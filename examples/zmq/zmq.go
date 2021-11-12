package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/libsv/go-bn/zmq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	z := zmq.NewNodeMQ(
		zmq.WithContext(ctx),
		zmq.WithHost("tcp://localhost:28332"),
		zmq.WithRaw(),
		zmq.WithErrorHandler(func(err error) {
			fmt.Println("OH NO", err)
		}),
	)

	if err := z.Subscribe(zmq.TopicHashTx, func(bb [][]byte) {
		fmt.Printf("tx hash: %s\n", hex.EncodeToString(bb[1]))
	}); err != nil {
		panic(err)
	}

	if err := z.Subscribe(zmq.TopicRawTx, func(bb [][]byte) {
		fmt.Printf("tx hex: %s\n", hex.EncodeToString(bb[1]))
	}); err != nil {
		panic(err)
	}

	for err := z.Connect(); err != nil; {
		time.Sleep(10 * time.Second)
		fmt.Println(err)
	}
}
