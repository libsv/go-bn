package zmq

import (
	"context"

	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
)

// MessageFunc a func in which the message is passed to.
type MessageFunc func(ctx context.Context, bb [][]byte)

// ErrorFunc a func in which an error is passed to.
type ErrorFunc func(ctx context.Context, err error)

// HashFunc a func in which `hashtx` and `hashblock` results are passed to.
type HashFunc func(ctx context.Context, hash string)

// DiscardFunc a func in which `hashtx` and `hashblock` results are passed to.
type DiscardFunc func(ctx context.Context, discard *MempoolDiscard)

// RawTxFunc a func in which `rawtx` results are parsed and passed to.
type RawTxFunc func(ctx context.Context, tx *bt.Tx)

// RawBlockFunc a func in which `rawblock` results are parsed and passed to.
type RawBlockFunc func(ctx context.Context, blk *bc.Block)

// MempoolDiscard a JSON representation of `discardfrommempool` and `removedfrommempoolblock`
// messages.
type MempoolDiscard struct {
	TxID         string `json:"txid"`
	Reason       string `json:"reason"`
	BlockHash    string `json:"blockhash"`
	CollidedWith struct {
		TxID string `json:"txid"`
		Size int    `json:"size"`
		Tx   *bt.Tx `json:"tx"`
	} `json:"collidedWith"`
}
