package models

import (
	"github.com/libsv/go-bn/models"
	"github.com/libsv/go-bt/v2"
)

type InternalFundRawTransaction struct {
	*models.FundRawTransaction
	Hex string `json:"hex"`
}

func (i *InternalFundRawTransaction) PostProcess() error {
	var err error
	i.Tx, err = bt.NewTxFromString(i.Hex)
	return err
}

type InternalSignRawTransaction struct {
	Hex string `json:"hex"`
	*models.SignedRawTransaction
}

func (i *InternalSignRawTransaction) PostProcess() error {
	var err error
	i.Tx, err = bt.NewTxFromString(i.Hex)
	return err
}
