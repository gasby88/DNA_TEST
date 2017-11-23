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
	DNA_RPC_GETIDENTITY         = "getidentity"
	DNA_RPC_GETIDENTITYCLAIM    = "getidentityclaim"
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

	DNA_SENDTRANSACTION     = "sendrawtransaction"
	DNA_HEARTBEAT           = "heartbeat"
	DNA_SMARTCONTRACTINVOKE = "InvokeTransaction"
)

const (
	DNA_ERR_OK                     = 0
	DNA_ERR_SESSION_EXPIRED        = 41001
	DNA_ERR_SERVICE_CEILING        = 41002
	DNA_ERR_ILLEGAL_DATAFORMAT     = 41003
	DNA_ERR_OAUTH_TIMEOUT          = 41004
	DNA_ERR_INVALID_METHOD         = 42001
	DNA_ERR_INVALID_PARAMS         = 42002
	DNA_ERR_INVALID_TOKEN          = 42003
	DNA_ERR_INVALID_TRANSACTION    = 43001
	DNA_ERR_INVALID_ASSET          = 43002
	DNA_ERR_INVALID_BLOCK          = 43003
	DNA_ERR_UNKNOWN_TRANSACTION    = 44001
	DNA_ERR_UNKNOWN_ASSET          = 44002
	DNA_ERR_UNKNOWN_BLOCK          = 44003
	DNA_ERR_INVALID_VERSION        = 45001
	DNA_ERR_INTERNAL_ERROR         = 45002
	DNA_ERR_OAUTH_INVALID_APPID    = 46001
	DNA_ERR_OAUTH_INVALID_CHECKVAL = 46002
	DNA_ERR_SMARTCODE_ERROR        = 47001

	ErrNoCode               = -2
	ErrNoError              = 0
	ErrUnknown              = -1
	ErrDuplicatedTx         = 1
	ErrDuplicateInput       = 45003
	ErrAssetPrecision       = 45004
	ErrTransactionBalance   = 45005
	ErrAttributeProgram     = 45006
	ErrTransactionContracts = 45007
	ErrTransactionPayload   = 45008
	ErrDoubleSpend          = 45009
	ErrTxHashDuplicate      = 45010
	ErrStateUpdaterVaild    = 45011
	ErrSummaryAsset         = 45012
	ErrXmitFail             = 45013
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
	DnaRpcIOError            = "internal IO error"
	DnaRpcAPIError           = "internal API error"
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
	DnaRpcIOError:            "",
	DnaRpcAPIError:           "",
}

func (this *DNAJsonRpcRes) HandleResult() ([]byte, error) {
	res := strings.Trim(string(this.Result), "\"")
	_, ok := DNARpcError[res]
	if ok {
		return nil, fmt.Errorf(res)
	}
	return []byte(res), nil
}
