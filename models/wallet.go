package models

import (
	"time"

	"github.com/libsv/go-bt/v2"
)

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
	LockingScript string    `json:"scriptPubKey"`
	Address       string    `json:"address"`
	Timestamp     time.Time `json:"timestamp"`
	RedeemScript  string    `json:"redeemscript"`
	PubKeys       []string  `json:"pubkeys"`
	Keys          []string  `json:"keys"`
	Internal      bool      `json:"internal"`
	WatchOnly     bool      `json:"watchonly"`
	Label         string    `json:"label"`
}

type OptsImportMulti struct {
	Rescan *bool `json:"rescan"`
}

func (o *OptsImportMulti) Args() []interface{} {
	if o.Rescan == nil {
		return []interface{}{false}
	}

	return []interface{}{*o.Rescan}
}

type ImportMulti struct {
	Success bool `json:"success"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type OptsImportPublicKey struct {
	Label  string
	Rescan *bool
}

func (o *OptsImportPublicKey) Args() []interface{} {
	aa := []interface{}{o.Label}
	if o.Rescan == nil {
		return aa
	}

	return append(aa, *o.Rescan)
}

type OptsKeypoolRefill struct {
	NewSize int
}

func (o *OptsKeypoolRefill) Args() []interface{} {
	if o.NewSize == 0 {
		return nil
	}

	return []interface{}{o.NewSize}
}

type OptsListAccounts struct {
	MinConf          int
	IncludeWatchOnly bool
}

func (o *OptsListAccounts) Args() []interface{} {
	aa := []interface{}{o.MinConf}
	if o.MinConf == 0 {
		aa[0] = 1
	}

	return append(aa, o.IncludeWatchOnly)
}

type LockUnspent struct {
	TxID string `json:"txid"`
	Vout int    `json:"vout"`
}

type ReceivedByAccount struct {
	InvolvesWatchOnly bool    `json:"involvesWatchOnly"`
	Account           string  `json:"account"`
	Amount            float64 `json:"amount"`
	Confirmations     int     `json:"confirmations"`
	Label             string  `json:"label"`
}

type OptsListReceivedBy struct {
	MinConf          int
	IncludeEmpty     bool
	IncludeWatchOnly bool
}

func (o *OptsListReceivedBy) Args() []interface{} {
	aa := []interface{}{o.MinConf}
	if o.MinConf == 0 {
		aa[0] = 1
	}

	return append(aa, o.IncludeEmpty, o.IncludeWatchOnly)
}

type ReceivedByAddress struct {
	InvolvesWatchOnly bool     `json:"involvesWatchOnly"`
	Address           string   `json:"address"`
	Account           string   `json:"account"`
	Amount            float64  `json:"amount"`
	Confirmations     int      `json:"confirmations"`
	Label             string   `json:"label"`
	TxIDs             []string `json:"txids"`
}

type SinceBlock struct {
	Txs []struct {
		Account       string  `json:"account"`
		Address       string  `json:"address"`
		Category      string  `json:"category"`
		Amount        float64 `json:"amount"`
		Generated     bool    `json:"generated"`
		Vout          int     `json:"vout"`
		Fee           float64 `json:"fee"`
		Confirmations int     `json:"confirmations"`
		BlockHash     string  `json:"blockhash"`
		BlockIndex    int     `json:"blockindex"`
		TxID          string  `json:"txid"`
		Time          uint64  `json:"time"`
		TimeReceived  uint64  `json:"timereceived"`
		Abandoned     bool    `json:"abandoned"`
		Comment       string  `json:"comment"`
		Label         string  `json:"label"`
		To            string  `json:"to"`
	} `json:"transactions"`
	LastBlock string `json:"lastblock"`
}

type OptsListSinceBlock struct {
	BlockHash           string
	TargetConfirmations int
	IncludeWatchOnly    bool
}

func (o *OptsListSinceBlock) Args() []interface{} {
	aa := []interface{}{o.BlockHash, o.TargetConfirmations}
	if o.TargetConfirmations == 0 {
		aa[1] = 1
	}

	return append(aa, o.IncludeWatchOnly)
}

type OptsListTransactions struct {
	Count            int
	Skip             int
	IncludeWatchOnly bool
}

func (o *OptsListTransactions) Args() []interface{} {
	count := o.Count
	if count == 0 {
		count = 10
	}

	return []interface{}{"*", o.Count, o.Skip, o.IncludeWatchOnly}
}

type OptsListUnspent struct {
	MinConf       int
	MaxConf       int
	Address       []string
	IncludeUnsafe *bool
}

func (o *OptsListUnspent) Args() []interface{} {
	aa := []interface{}{o.MinConf, o.MaxConf}
	if o.MinConf == 0 {
		o.MinConf = 1
	}
	if o.MaxConf == 0 {
		o.MaxConf = 9999999
	}

	if o.Address != nil && len(o.Address) > 0 {
		aa = append(aa, o.Address)
	}

	if o.IncludeUnsafe == nil {
		return aa
	}

	if len(aa) == 2 {
		aa = append(aa, []string{})
	}

	return append(aa, o.IncludeUnsafe)
}

type OptsLockUnspent struct {
	Txs []LockUnspent
}

func (o *OptsLockUnspent) Args() []interface{} {
	if o.Txs == nil || len(o.Txs) == 0 {
		return nil
	}

	return []interface{}{o.Txs}
}

type OptsMove struct {
	Comment string
}

func (o *OptsMove) Args() []interface{} {
	if o.Comment != "" {
		return []interface{}{"", o.Comment}
	}

	return nil
}

type OptsSendFrom struct {
	MinConf   int
	Comment   string
	CommentTo string
}

func (o *OptsSendFrom) Args() []interface{} {
	aa := []interface{}{o.MinConf}
	if aa[0] == 0 {
		aa[0] = 1
	}

	return append(aa, o.Comment, o.CommentTo)
}

type OptsSendMany struct {
	MinConf         int
	Comment         string
	SubtractFeeFrom []string
}

func (o *OptsSendMany) Args() []interface{} {
	aa := []interface{}{o.MinConf, o.Comment}
	if o.MinConf == 0 {
		aa[0] = 1
	}

	if o.SubtractFeeFrom != nil && len(o.SubtractFeeFrom) > 1 {
		aa = append(aa, o.SubtractFeeFrom)
	}

	return aa
}

type OptsSendToAddress struct {
	Comment         string
	CommentTo       string
	SubtractFeeFrom []string
}

func (o *OptsSendToAddress) Args() []interface{} {
	aa := []interface{}{o.Comment, o.CommentTo}

	if o.SubtractFeeFrom != nil && len(o.SubtractFeeFrom) > 1 {
		aa = append(aa, o.SubtractFeeFrom)
	}

	return aa
}
