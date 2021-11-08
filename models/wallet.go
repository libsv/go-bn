package models

import "github.com/libsv/go-bt/v2"

type DumpWallet struct {
	FileName string `json:"filename"`
}

type OptsBalance struct {
	Account              string
	MinimumConfirmations uint32
	IncludeWatchOnly     bool
}

func (o *OptsBalance) Args() []interface{} {
	aa := []interface{}{o.Account, o.MinimumConfirmations, o.IncludeWatchOnly}
	if o.MinimumConfirmations == 0 {
		aa[1] = 1
	}

	return aa
}

type OptsNewAddress struct {
	Account string
}

func (o *OptsNewAddress) Args() []interface{} {
	return []interface{}{o.Account}
}

type Transaction struct {
	Amount          float64       `json:"amount"`
	Fee             float64       `json:"fee"`
	Confirmations   uint32        `json:"confirmations"`
	BlockHash       string        `json:"blockhash"`
	BlockIndex      uint32        `json:"blockindex"`
	BlockTime       uint64        `json:"blocktime"`
	TxID            string        `json:"txid"`
	WalletConflicts []interface{} `json:"walletconflicts"`
	Time            uint64        `json:"time"`
	TimeReceived    uint64        `json:"timereceived"`
	Details         []struct {
		Account   string  `json:"account"`
		Address   string  `json:"address"`
		Category  string  `json:"category"`
		Amount    float64 `json:"amount"`
		Label     string  `json:"label"`
		Vout      uint32  `json:"vout"`
		Fee       float64 `json:"fee"`
		Abandoned bool    `json:"abandoned"`
	} `json:"details"`
	Tx *bt.Tx `json:"tx"`
}

type WalletInfo struct {
	WalletName            string  `json:"walletname"`
	WalletVersion         uint64  `json:"walletversion"`
	Balance               float64 `json:"balance"`
	UnconfirmedBalance    float64 `json:"unconfirmed_balance"`
	ImmatureBalance       float64 `json:"immature_balance"`
	TxCount               uint64  `json:"txcount"`
	KeypoolOldest         uint64  `json:"keypoololdest"`
	KeypoolSize           uint64  `json:"keypoolsize"`
	KeypoolSizeHDInternal uint32  `json:"keypoolsize_hd_internal"`
	PayTxFee              float64 `json:"paytxfee"`
	HDMasterKeyID         string  `json:"hdmasterkeyid"`
}

type OptsImportAddress struct {
	Label  string
	Rescan *bool
	P2SH   bool
}

func (o *OptsImportAddress) Args() []interface{} {
	aa := []interface{}{o.Label, true, o.P2SH}
	if o.Rescan != nil {
		aa[1] = o.Rescan
	}

	return aa
}

type OptsImportPrivateKey struct {
	Label  string
	Rescan *bool
}

func (o *OptsImportPrivateKey) Args() []interface{} {
	aa := []interface{}{o.Label}
	if o.Rescan != nil {
		aa = append(aa, o.Rescan)
	}

	return aa
}

type ImportMultiRequest struct {
	LockingScript string
}
