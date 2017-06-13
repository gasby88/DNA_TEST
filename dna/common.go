package dna

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	DNA_RPC_GETVERSION          = "getversion"
	DNA_RPC_GETTRANSACTION      = "getrawtransaction"
	DNA_RPC_SENDTRANSACTION     = "sendrawtransaction"
	DNA_RPC_GETBLOCK            = "getblock"
	DNA_RPC_GETBLOCKCOUNT       = "getblockcount"
	DNA_RPC_GETBLOCKHASH        = "getblockhash"
	DNA_RPC_GETUNSPENDOUTPUT    = "getunspendoutput"
	DNA_RPC_GETCURRENTBLOCKHASH = "getbestblockhash"
)

type DNAJsonRpcRes struct {
	Id      string          `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
}

var DNARpcError map[string]string = map[string]string{
	"invalid hash":        "",
	"invalid block":       "",
	"invalid transaction": "",
	"invalid parameter":   "",
	"unknown block":       "",
	"unknown transaction": "",
	"Unsupported":         "",
	"internal error":      "",
	"null":                "",
}

func (this *DNAJsonRpcRes) HandleResult() ([]byte, error) {
	res := strings.Trim(string(this.Result), "\"")
	_, ok := DNARpcError[res]
	if ok {
		return nil, fmt.Errorf(res)
	}
	return []byte(res), nil
}
