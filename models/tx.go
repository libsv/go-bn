package models

import (
	"encoding/json"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"
)

type Output struct {
	BestBlock     string
	Confirmations uint32
	Coinbase      bool

	*bt.Output
}

func (o *Output) NodeJSON() interface{} {
	return o
}

func (o *Output) UnmarshalJSON(b []byte) error {
	oj := struct {
		BestBlock     string `json:"bestblock"`
		Confirmations uint32 `json:"confirmations"`
		Coinbase      bool   `json:"coinbase"`
	}{}

	if err := json.Unmarshal(b, &oj); err != nil {
		return err
	}

	var out bt.Output
	if err := json.Unmarshal(b, out.NodeJSON()); err != nil {
		return err
	}

	o.BestBlock = oj.BestBlock
	o.Confirmations = oj.Confirmations
	o.Coinbase = oj.Coinbase
	*o.Output = out

	return nil
}

type OutputSetInfo struct {
	Height         uint32  `json:"height"`
	BestBlock      string  `json:"bestblock"`
	Transactions   uint32  `json:"transactions"`
	OutputCount    uint32  `json:"txouts"`
	BogoSize       uint32  `json:"bogosize"`
	HashSerialised string  `json:"hash_serialized"`
	DiskSize       uint32  `json:"disk_size"`
	TotalAmount    float64 `json:"total_amount"`
}

type OptsOutput struct {
	IncludeMempool bool
}

func (o *OptsOutput) Args() []interface{} {
	return []interface{}{o.IncludeMempool}
}

type ParamsCreateRawTransaction struct {
	Outputs []*bt.Output
	mainnet bool
}

func (p *ParamsCreateRawTransaction) Args() []interface{} {
	outputs := make(map[string]float64, len(p.Outputs))
	for _, o := range p.Outputs {
		pkh, err := o.LockingScript.PublicKeyHash()
		if err != nil {
			outputs["invalid locking script"] = float64(o.Satoshis) / 100000000
			continue
		}
		addr, err := bscript.NewAddressFromPublicKeyHash(pkh, p.mainnet)
		if err != nil {
			outputs["invalid locking script"] = float64(o.Satoshis) / 100000000
		}
		outputs[addr.AddressString] = float64(o.Satoshis) / 100000000
	}

	return []interface{}{outputs}
}

func (p *ParamsCreateRawTransaction) SetIsMainnet(b bool) {
	p.mainnet = b
}

type FundRawTransaction struct {
	Fee            uint64 `json:"fee"`
	ChangePosition int    `json:"changeposition"`
	Tx             *bt.Tx
}

type OptsFundRawTransaction struct {
	ChangeAddress          string   `json:"changeAddress,omitempty"`
	ChangePosition         int      `json:"changePosition,omitempty"`
	IncludeWatching        bool     `json:"includeWatching,omitempty"`
	LockUnspents           bool     `json:"lockUnspents,omitempty"`
	ReserveChangeKey       *bool    `json:"reserveChangeKey,omitempty"`
	FeeRate                uint64   `json:"feeRate,omitempty"`
	SubtractFeeFromOutputs []uint64 `json:"subtractFeeFromOutputs,omitempty"`
}

func (o *OptsFundRawTransaction) Args() []interface{} {
	return []interface{}{o}
}

type SendRawTransaction struct {
	Hex string
	Tx  *bt.Tx
}

func (s *SendRawTransaction) PostProcess() error {
	var err error
	s.Tx, err = bt.NewTxFromString(s.Hex)
	return err
}

type SignedRawTransaction struct {
	Tx       *bt.Tx `json:"tx"`
	Complete bool   `json:"complete"`
	Errors   []struct {
		TxID            string `json:"txid"`
		Vout            int    `json:"vout"`
		UnlockingScript string `json:"scriptSig"`
		Sequence        uint32 `json:"sequence"`
		Error           string `json:"error"`
	} `json:"errors"`
}

type OptsSignRawTransaction struct {
	From        bt.UTXOs
	PrivateKeys []string
	SigHashType sighash.Flag
}

func (o *OptsSignRawTransaction) Args() []interface{} {
	aa := []interface{}{[]interface{}{}, []interface{}{}}
	if o.From != nil && len(o.From) > 0 {
		aa[0] = o.From.NodeJSON()
	}
	if o.PrivateKeys != nil && len(o.PrivateKeys) > 0 {
		aa[0] = o.PrivateKeys
	}
	return append(aa, o.SigHashType.String())
}

type OptsSendRawTransaction struct {
	AllowHighFees bool
	CheckFee      bool
}

func (o *OptsSendRawTransaction) Args() []interface{} {
	return []interface{}{o.AllowHighFees, !o.CheckFee}
}
