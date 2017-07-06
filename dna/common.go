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

const (
	DNA_API_GETCONNCOUNT     = "/api/v1/node/connectioncount"
	DNA_API_GETBLOCKBYHEIGHT = "/api/v1/block/details/height"
	DNA_API_GETBLOCKBYHASH   = "/api/v1/block/details/hash"
	DNA_API_GETBLOCKCOUNT    = "/api/v1/block/height"
	DNA_API_GETTRANSACTION   = "/api/v1/transaction"
	DNA_API_GETASSET         = "/api/v1/asset"
	DNA_API_SENDTRANSACTION  = "/api/v1/transaction"

	//Api_Getconnectioncount = "/api/v1/node/connectioncount"
	//Api_Getblockbyheight   = "/api/v1/block/details/height/:height"
	//Api_Getblockbyhash     = "/api/v1/block/details/hash/:hash"
	//Api_Getblockheight     = "/api/v1/block/height"
	//Api_Gettransaction     = "/api/v1/transaction/:hash"
	//Api_Getasset           = "/api/v1/asset/:hash"
	//Api_Restart            = "/api/v1/restart"
	//Api_SendRawTransaction = "/api/v1/transaction"
	//Api_OauthServerAddr    = "/api/v1/config/oauthserver/addr"
	//Api_NoticeServerAddr   = "/api/v1/config/noticeserver/addr"
	//Api_NoticeServerState  = "/api/v1/config/noticeserver/state"
)

type DNAJsonRpcRes struct {
	Id      string          `json:"id"`
	JsonRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
}

const (
	DnaRpcInvalidHash        = "invalid hash"
	DnaRpcInvalidBlock       = "invalid block"
	DnaRpcInvalidTransaction = "invalid transaction"
	DnaRpcInvalidParameter   = "invalid parameter"
	DnaRpcUnknownBlock       = "unknown block"
	DnaRpcUnknownTransaction = "unknown transaction"
	DnaRpcNil                = "null"
	DnaRpcUnsupported        = "Unsupported"
	DnaRpcInternalError      = "internal error"
)

var DNARpcError map[string]string = map[string]string{
	DnaRpcInvalidHash:        "",
	DnaRpcInvalidBlock:       "",
	DnaRpcInvalidTransaction: "",
	DnaRpcInvalidParameter:   "",
	DnaRpcUnknownBlock:       "",
	DnaRpcUnknownTransaction: "",
	DnaRpcUnsupported:        "",
	DnaRpcInternalError:      "",
	DnaRpcNil:                "",
}

func (this *DNAJsonRpcRes) HandleResult() ([]byte, error) {
	res := strings.Trim(string(this.Result), "\"")
	_, ok := DNARpcError[res]
	if ok {
		return nil, fmt.Errorf(res)
	}
	return []byte(res), nil
}
