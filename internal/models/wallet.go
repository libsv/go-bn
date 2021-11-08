package models

import (
	"encoding/json"

	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bn/models"
	"github.com/libsv/go-bt/v2"
)

type InternalDumpPrivateKey struct {
	WIF *wif.WIF
}

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

type InternalTransaction struct {
	*models.Transaction
	Hex string `json:"hex"`
}

func (i *InternalTransaction) PostProcess() error {
	var err error
	i.Tx, err = bt.NewTxFromString(i.Hex)
	return err
}
