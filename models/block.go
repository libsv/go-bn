package models

import (
	"encoding/json"

	"github.com/libsv/go-bt/v2"
)

// BlockDecodeHeader model.
type BlockDecodeHeader struct {
	Txs []string `json:"tx"`
	BlockHeader
}

// Block model.
type Block struct {
	Txs bt.Txs `json:"tx"`
	BlockHeader
}

// UnmarshalJSON unmarshal response.
func (b *Block) UnmarshalJSON(bb []byte) error {
	bj := struct {
		Txs json.RawMessage `json:"tx"`
		BlockHeader
	}{}
	if err := json.Unmarshal(bb, &bj); err != nil {
		return err
	}

	var txs bt.Txs
	if err := json.Unmarshal(bj.Txs, txs.NodeJSON()); err != nil {
		return err
	}

	b.Txs = txs
	b.BlockHeader = bj.BlockHeader
	return nil
}

// BlockHeader model.
type BlockHeader struct {
	Hash              string  `json:"hash"`
	Confirmations     uint64  `json:"confirmations"`
	Size              uint64  `json:"size"`
	Height            uint64  `json:"height"`
	Version           uint64  `json:"version"`
	VersionHex        string  `json:"versionHex"`
	NumTx             uint64  `json:"num_tx"`
	Time              uint64  `json:"time"`
	MedianTime        uint64  `json:"mediantime"`
	Nonce             uint64  `json:"nonce"`
	Bits              string  `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
	Chainwork         string  `json:"chainwork"`
	PreviousBlockHash string  `json:"previousblockhash"`
	NextBlockHash     string  `json:"nextblockhash"`
}

// BlockTemplate model.
type BlockTemplate struct {
	Capabilities      []string `json:"capabilities"`
	Version           uint64   `json:"version"`
	PreviousBlockHash string   `json:"previousblockhash"`
	Transactions      []string `json:"transactions"`
	CoinbaseAux       struct {
		Flags string `json:"flags"`
	} `json:"coinbaseaux"`
	CoinbaseValue uint64   `json:"coinbasevalue"`
	LongPollID    string   `json:"longpollid"`
	Target        string   `json:"target"`
	MinTime       uint64   `json:"mintime"`
	Mutable       []string `json:"mutable"`
	NonceRange    string   `json:"noncerange"`
	SizeLimit     uint64   `json:"sizelimit"`
	CurTime       uint64   `json:"curtime"`
	Bits          string   `json:"bits"`
	Height        uint64   `json:"height"`
}

// BlockTemplateRequest model.
type BlockTemplateRequest struct {
	Mode         string
	Capabilities []string
}

// Args convert struct into optional positional arguments.
func (r *BlockTemplateRequest) Args() []interface{} {
	return []interface{}{r}
}
