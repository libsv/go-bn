package models

import (
	"encoding/json"

	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bn/models"
	"github.com/libsv/go-bt/v2"
)

// InternalDumpPrivateKey the true to form dumpprivkey response from the bitcoin node.
type InternalDumpPrivateKey struct {
	WIF *wif.WIF
}

// UnmarshalJSON unmarshal the response.
func (i *InternalDumpPrivateKey) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	w, err := wif.DecodeWIF(s)
	if err != nil {
		return err
	}

	i.WIF = w
	return nil
}

// InternalTransaction the true to form transaction response from the bitcoin node.
type InternalTransaction struct {
	*models.Transaction
	Hex string `json:"hex"`
}

// PostProcess an RPC response.
func (i *InternalTransaction) PostProcess() error {
	var err error
	i.Tx, err = bt.NewTxFromString(i.Hex)
	return err
}
