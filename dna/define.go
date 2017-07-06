package dna

import (
	. "DNA/common"
	. "DNA/core/asset"
	"DNA/core/transaction"
	"DNA/net/httpjsonrpc"
	"encoding/json"
)

type BlockInfo struct {
	Hash         string
	BlockData    *BlockHead
	Transactions []*Transactions
}

type BlockHead struct {
	Version          uint32
	PrevBlockHash    string
	TransactionsRoot string
	Timestamp        uint32
	Height           uint32
	ConsensusData    uint64
	NextBookKeeper   string
	Program          ProgramInfo
	Hash             string
}

type TxAttributeInfo struct {
	Usage byte
	Data  string
}

type UTXOTxInputInfo struct {
	ReferTxID          string
	ReferTxOutputIndex uint16
}

type BalanceTxInputInfo struct {
	AssetID     string
	Value       Fixed64
	ProgramHash string
}

type TxoutputInfo struct {
	AssetID     string
	Value       Fixed64
	ProgramHash string
}

type ProgramInfo struct {
	Code      string
	Parameter string
}

type TxoutputMap struct {
	Key   Uint256
	Txout []TxoutputInfo
}

type AmountMap struct {
	Key   Uint256
	Value Fixed64
}

type Transactions struct {
	TxType            transaction.TransactionType
	PayloadVersion    byte
	Payload           json.RawMessage
	Attributes        []TxAttributeInfo
	UTXOInputs        []UTXOTxInputInfo
	BalanceInputs     []BalanceTxInputInfo
	Outputs           []TxoutputInfo
	Programs          []ProgramInfo
	AssetOutputs      []TxoutputMap
	AssetInputAmount  []AmountMap
	AssetOutputAmount []AmountMap
	Hash              string
}

type PayloadRegisterAssetInfo struct {
	Asset      *Asset
	Amount     Fixed64
	Issuer     httpjsonrpc.IssuerInfo
	Controller string
}

type PayloadRecord struct {
	RecordType string
	RecordData string
}

type PayloadDeployCodeInfo struct {
	Code        *httpjsonrpc.FunctionCodeInfo
	Name        string
	CodeVersion string
	Author      string
	Email       string
	Description string
}

type UnspendUTXO struct {
	ReferTxID          Uint256
	ReferTxOutputIndex uint16
	AssetID            Uint256
	Value              Fixed64
	ProgramHash        Uint160
}
