package models

import (
	"fmt"
)

type blockVerbosity string

const (
	VerbosityRawBlock                blockVerbosity = "RAW_BLOCK"
	VerbosityDecodeHeader            blockVerbosity = "DECODE_HEADER"
	VerbosityDecodeTransactions      blockVerbosity = "DECODE_TRANSACTIONS"
	VerbosityDecodeHeaderAndCoinbase blockVerbosity = "DECODE_HEADER_AND_COINBASE"
)

type merkleProofTargetType string

const (
	MerkleProofTargetTypeHash       merkleProofTargetType = "hash"
	MerkleProofTargetTypeHeader     merkleProofTargetType = "header"
	MerkleProofTargetTypeMerkleRoot merkleProofTargetType = "merkleroot"
)

type Request struct {
	ID      string        `json:"id"`
	JSONRpc string        `json:"jsonRpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

type Response struct {
	Result interface{} `json:"result"`
	Error  *Error      `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

type OptsChainTxStats struct {
	NumBlocks uint32
	BlockHash string
}

func (o *OptsChainTxStats) Args() []interface{} {
	aa := []interface{}{o.NumBlocks}
	if o.BlockHash != "" {
		aa = append(aa, o.BlockHash)
	}
	return aa
}

type OptsMerkleProof struct {
	FullTx     bool
	TargetType merkleProofTargetType
}

func (o *OptsMerkleProof) Args() []interface{} {
	aa := []interface{}{o.FullTx}
	if o.TargetType != "" {
		aa = append(aa, o.TargetType)
	}

	return aa
}

type OptsLegacyMerkleProof struct {
	BlockHash string
}

func (o *OptsLegacyMerkleProof) Args() []interface{} {
	return []interface{}{o.BlockHash}
}

type OptsGenerate struct {
	MaxTries uint32
}

func (o *OptsGenerate) Args() []interface{} {
	return []interface{}{o.MaxTries}
}
