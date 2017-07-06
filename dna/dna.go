package dna

import (
	"DNA/account"
	. "DNA/common"
	"DNA/core/asset"
	"DNA/core/contract"
	"DNA/core/ledger"
	"DNA/core/signature"
	"DNA/core/transaction"
	"DNA/core/transaction/payload"
	"DNA/crypto"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	//log4 "github.com/alecthomas/log4go"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	crypto.SetAlg("P256R1")
}

type Dna struct {
	qid          uint64
	rpcAddresses []string
	client       *http.Client
}

func NewDna(rpcAddresses []string) *Dna {
	return &Dna{
		rpcAddresses: rpcAddresses,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   50,
				DisableKeepAlives:     false, //启动keepalive
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
			},
			Timeout: time.Second * 300,
		},
	}
}

func (this *Dna) GetVersion() (string, error) {
	data, err := this.sendRpcRequest(DNA_RPC_GETVERSION, []interface{}{})
	if err != nil {
		return "", fmt.Errorf("SendRpcRequest error:%s", err)
	}
	return string(data), nil
}

func (this *Dna) CreateAsset(
	name string,
	precision byte,
	assetType asset.AssetType,
	recordType asset.AssetRecordType) *asset.Asset {
	return &asset.Asset{
		Name:       name,
		Precision:  precision,
		AssetType:  assetType,
		RecordType: recordType,
	}
}

func (this *Dna) GetBlockByHash(hash Uint256) (*ledger.Block, error) {
	blockHash := Uint256ToString(hash)
	data, err := this.sendRpcRequest(DNA_RPC_GETBLOCK, []interface{}{Uint256ToString(hash)})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	blockInfo := &BlockInfo{}
	err = json.Unmarshal(data, blockInfo)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal BlockInfo:%s error:%s", blockInfo, err)
	}
	block, err := ParseBlock(blockInfo)
	if err != nil {
		return nil, fmt.Errorf("ParseBlock Hash:%x error:%s", blockHash, err)
	}
	return block, nil
}

func (this *Dna) GetBlockByHeight(height uint32) (*ledger.Block, error) {
	data, err := this.sendRpcRequest(DNA_RPC_GETBLOCK, []interface{}{height})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	blockInfo := &BlockInfo{}
	err = json.Unmarshal(data, blockInfo)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal BlockInfo:%s error:%s", blockInfo, err)
	}
	block, err := ParseBlock(blockInfo)
	if err != nil {
		return nil, fmt.Errorf("ParseBlock Hright:%v error:%s", height, err)
	}
	return block, nil
}

func (this *Dna) GetBlockHash(height uint32) (Uint256, error) {
	data, err := this.sendRpcRequest(DNA_RPC_GETBLOCKHASH, []interface{}{height})
	if err != nil {
		return Uint256{}, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	hash, err := ParseUint256FromString(string(data))
	if err != nil {
		return Uint256{}, fmt.Errorf("ParseUint256FromString Hash:%s error:%s", data, err)
	}
	return hash, nil
}

func (this *Dna) GetCurrentBlockHash() (Uint256, error) {
	data, err := this.sendRpcRequest(DNA_RPC_GETCURRENTBLOCKHASH, []interface{}{})
	if err != nil {
		return Uint256{}, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	hash, err := ParseUint256FromString(string(data))
	if err != nil {
		return Uint256{}, fmt.Errorf("ParseUint256FromString:%s error:%s", hash, err)
	}
	return hash, nil
}

func (this *Dna) GetBlockCount() (uint32, error) {
	data, err := this.sendRpcRequest(DNA_RPC_GETBLOCKCOUNT, []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	count := uint32(0)
	err = json.Unmarshal(data, &count)
	if err != nil {
		return 0, fmt.Errorf("json.Unmarshal Count:%s error:%s", data, err)
	}
	return count, nil
}

func (this *Dna) NewAssetRegisterTransaction(asset *asset.Asset,
	amount Fixed64,
	issuer,
	controllerAccount *account.Account) (*transaction.Transaction, error) {
	controller, err := contract.CreateSignatureContract(controllerAccount.PubKey())
	if err != nil {
		return nil, fmt.Errorf("CreateSignatureContract error:%s", err)
	}
	tx, err := transaction.NewRegisterAssetTransaction(asset, amount, issuer.PubKey(), controller.ProgramHash)
	if err != nil {
		return nil, fmt.Errorf("NewRegisterAssetTransaction error:%s", err)
	}
	this.setNonce(tx)
	return tx, nil
}

func (this *Dna) NewIssueAssetTransaction(txOutputs []*transaction.TxOutput) (*transaction.Transaction, error) {
	tx, err := transaction.NewIssueAssetTransaction(txOutputs)
	if err != nil {
		return nil, fmt.Errorf("NewIssueAssetTransaction error:%s", err)
	}
	this.setNonce(tx)
	return tx, nil
}

func (this *Dna) NewTransferAssetTransaction(inputs []*transaction.UTXOTxInput,
	outputs []*transaction.TxOutput) (*transaction.Transaction, error) {
	tx, err := transaction.NewTransferAssetTransaction(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("NewTransferAssetTransaction error:%s", err)
	}
	this.setNonce(tx)
	return tx, nil
}

func (this *Dna) NewRecordTransaction(recordType string, recordData []byte) (*transaction.Transaction, error) {
	tx, err := transaction.NewRecordTransaction(recordType, recordData)
	if err != nil {
		return nil, fmt.Errorf("NewRecordTransaction error:%s", err)
	}
	this.setNonce(tx)
	return tx, nil
}

func (this *Dna) setNonce(tx *transaction.Transaction) {
	attr := transaction.NewTxAttribute(transaction.Nonce, []byte(fmt.Sprintf("%d", rand.Int63())))
	tx.Attributes = append(tx.Attributes, &attr)
}

func (this *Dna) SendTransaction(account *account.Account, tx *transaction.Transaction) (Uint256, error) {
	err := this.SignTransaction(account, tx)
	if err != nil {
		return Uint256{}, fmt.Errorf("SignTransaction error:%s", err)
	}

	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return Uint256{}, fmt.Errorf("Serialize error:%s", err)
	}

	txData := hex.EncodeToString(buffer.Bytes())
	data, err := this.sendRpcRequest(DNA_RPC_SENDTRANSACTION, []interface{}{txData})
	if err != nil {
		return Uint256{}, err
	}

	hash, err := ParseUint256FromString(string(data))
	if err != nil {
		return Uint256{}, fmt.Errorf("ParseUint256FromString Hash:%s error:%s", data, err)
	}
	return hash, nil
}

func (this *Dna) SignTransaction(signer *account.Account, tx *transaction.Transaction) error {
	signature, err := signature.SignBySigner(tx, signer)
	if err != nil {
		return fmt.Errorf("SignBySigner error:%s", err)
	}
	transactionContract, err := contract.CreateSignatureContract(signer.PubKey())
	if err != nil {
		return fmt.Errorf("CreateSignatureContract error:%s", err)
	}
	programHashes, err := this.GetTransactionProgramHashes(tx)
	if err != nil {
		return fmt.Errorf("GetTransactionProgramHashes error:%s", err)
	}
	ctx, err := this.NewContractContext(tx, programHashes)
	if err != nil {
		return fmt.Errorf("NewContractContext error:%s", err)
	}
	err = ctx.AddContract(transactionContract, signer.PubKey(), signature)
	if err != nil {
		return fmt.Errorf("AddContract error:%s", err)
	}
	tx.SetPrograms(ctx.GetPrograms())
	return nil
}

func (this *Dna) GetTransactionProgramHashes(tx *transaction.Transaction) ([]Uint160, error) {
	hashs := []Uint160{}
	uniqHashes := []Uint160{}
	// add inputUTXO's transaction
	referenceWithUTXO_Output, err := this.GetTransactionReference(tx)
	if err != nil {
		return nil, fmt.Errorf("Transction GetReference error:%s", err)
	}
	for _, output := range referenceWithUTXO_Output {
		programHash := output.ProgramHash
		hashs = append(hashs, programHash)
	}
	for _, attribute := range tx.Attributes {
		if attribute.Usage != transaction.Script {
			continue
		}
		dataHash, err := Uint160ParseFromBytes(attribute.Data)
		if err != nil {
			return nil, fmt.Errorf("Uint160ParseFromBytes error:%s", err)
		}
		hashs = append(hashs, Uint160(dataHash))
	}
	switch tx.TxType {
	case transaction.RegisterAsset:
		issuer := tx.Payload.(*payload.RegisterAsset).Issuer
		signatureRedeemScript, err := contract.CreateSignatureRedeemScript(issuer)
		if err != nil {
			return nil, fmt.Errorf("CreateSignatureRedeemScript error:%s", err)
		}
		astHash, err := ToCodeHash(signatureRedeemScript)
		if err != nil {
			return nil, fmt.Errorf("ToCodeHash error:%s", err)
		}
		hashs = append(hashs, astHash)
	case transaction.IssueAsset:
		result := tx.GetMergedAssetIDValueFromOutputs()
		if err != nil {
			return nil, fmt.Errorf("GetMergedAssetIDValueFromOutputs error:%s", err)
		}
		for k, _ := range result {
			regTx, err := this.GetTransaction(k)
			if err != nil {
				return nil, fmt.Errorf("GetTransaction TxHash:%x error:%s", k, err)
			}
			if regTx.TxType != transaction.RegisterAsset {
				return nil, errors.New("Transaction is not RegisterAsset")
			}

			regPayload := regTx.Payload.(*payload.RegisterAsset)
			hashs = append(hashs, regPayload.Controller)
		}
	case transaction.TransferAsset:
	case transaction.Record:
	case transaction.BookKeeper:
	default:
	}
	//remove dupilicated hashes
	uniq := make(map[Uint160]bool)
	for _, v := range hashs {
		uniq[v] = true
	}
	for k, _ := range uniq {
		uniqHashes = append(uniqHashes, k)
	}
	sort.Sort(ByProgramHashes(uniqHashes))
	return uniqHashes, nil
}

func (this *Dna) NewContractContext(data signature.SignableData, programHashes ...[]Uint160) (*contract.ContractContext, error) {
	var proHashes []Uint160
	var err error
	if len(programHashes) > 0 {
		proHashes = programHashes[0]
	} else {
		proHashes, err = data.GetProgramHashes()
		if err != nil {
			return nil, fmt.Errorf("GetProgramHashes error:%s", err)
		}
	}
	hashLen := len(proHashes)
	return &contract.ContractContext{
		Data:            data,
		ProgramHashes:   proHashes,
		Codes:           make([][]byte, hashLen),
		Parameters:      make([][][]byte, hashLen),
		MultiPubkeyPara: make([][]contract.PubkeyParameter, hashLen),
	}, nil
}

func (this *Dna) GetTransaction(txHash Uint256) (*transaction.Transaction, error) {
	data, err := this.sendRpcRequest(DNA_RPC_GETTRANSACTION, []interface{}{Uint256ToString(txHash)})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	txStr := &Transactions{}
	err = json.Unmarshal(data, txStr)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal Transactions:%s error:%s", data, err)
	}
	tx, err := ParseTransaction(txStr)
	if err != nil {
		return nil, fmt.Errorf("ParseTransaction:%+v error:%s", txStr, err)
	}
	return tx, nil
}

func (this *Dna) GetUnspendOutput(assetHash Uint256, programHash Uint160) ([]*UnspendUTXO, error) {
	data, err := this.sendRpcRequest(DNA_RPC_GETUNSPENDOUTPUT, []interface{}{Uint160ToString(programHash), Uint256ToString(assetHash)})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	if string(data) == "{}" {
		return nil, nil
	}
	outputMap := make(map[string]json.RawMessage, 0)
	err = json.Unmarshal(data, &outputMap)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal map[string]json.RawMessage:%s error:%s", data, err)
	}
	unspents := make([]*UnspendUTXO, 0, len(outputMap))
	for k, o := range outputMap {
		output := &TxoutputInfo{}
		err := json.Unmarshal(o, output)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal TxoutputInfo:%s error:%s", output, err)
		}
		txOutput, err := ParseTransactionOutputs(output)
		ks := strings.Split(k, ":")
		if len(ks) != 2 {
			return nil, fmt.Errorf("UnspentUTXO key:%s error", k)
		}
		referId, err := ParseUint256FromString(ks[0])
		if err != nil {
			return nil, fmt.Errorf("ParseUint256FromString:%x error:%s", ks[0], err)
		}
		index, err := strconv.ParseInt(ks[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("strconv.ParseInt:%s error:%s", ks[1], err)
		}
		unspent := &UnspendUTXO{
			AssetID:            txOutput.AssetID,
			Value:              txOutput.Value,
			ProgramHash:        txOutput.ProgramHash,
			ReferTxID:          referId,
			ReferTxOutputIndex: uint16(index),
		}
		if err != nil {
			return nil, fmt.Errorf("ParseTransactionOutputs:%s error:%s", txOutput, err)
		}
		unspents = append(unspents, unspent)
	}
	return unspents, nil
}

func (this *Dna) WaitForGenerateBlock(timeout time.Duration, blockCount ...uint32) (bool, error) {
	count := uint32(2)
	if len(blockCount) > 0 && blockCount[0] > 0 {
		count = blockCount[0]
	}
	blockHeight, err := this.GetBlockCount()
	if err != nil {
		return false, fmt.Errorf("GetBlockCount error:%s", err)
	}
	secs := int(timeout / time.Second)
	if secs <= 0 {
		secs = 1
	}
	ok := false
	for i := 0; i < secs; i++ {
		time.Sleep(time.Second)
		curBlockHeigh, err := this.GetBlockCount()
		if err != nil {
			continue
		}
		if curBlockHeigh-blockHeight >= count {
			ok = true
			break
		}
	}
	return ok, nil
}

func (this *Dna) MakeAssetAmount(rawAmont float64) Fixed64 {
	return Fixed64(rawAmont * 100000000)
}

func (this *Dna) GetRawAssetAmount(assetAmount Fixed64) float64 {
	return float64(assetAmount) / 100000000
}

func (this *Dna) GetAccountProgramHash(account *account.Account) (Uint160, error) {
	ctr, err := contract.CreateSignatureContract(account.PubKey())
	if err != nil {
		return Uint160{}, fmt.Errorf("CreateSignatureContract error:%s", err)
	}
	return ctr.ProgramHash, nil
}

func (this *Dna) getQid() string {
	return fmt.Sprintf("%d", atomic.AddUint64(&this.qid, 1))
}

func (this *Dna) getRpcAddress() string {
	if len(this.rpcAddresses) == 0 {
		return ""
	}
	return this.rpcAddresses[0]
}

func (this *Dna) GetTransactionReference(tx *transaction.Transaction) (map[*transaction.UTXOTxInput]*transaction.TxOutput, error) {
	if tx.TxType == transaction.RegisterAsset {
		return nil, nil
	}
	//UTXO input /  Outputs
	reference := make(map[*transaction.UTXOTxInput]*transaction.TxOutput)
	// Key index，v UTXOInput
	for _, utxo := range tx.UTXOInputs {
		referTx, err := this.GetTransaction(utxo.ReferTxID)
		if err != nil {
			return nil, fmt.Errorf("GetTransaction refer txHash:%x", utxo.ReferTxID)
		}
		index := utxo.ReferTxOutputIndex
		reference[utxo] = referTx.Outputs[index]
	}
	return reference, nil

}
func (this *Dna) sendRpcRequest(method string, params []interface{}) ([]byte, error) {
	data, err := this.Call(this.getRpcAddress(), method, this.getQid(), params)
	//if method == DNA_RPC_SENDTRANSACTION {
	//	log4.Debug("Call:%s params:%+v", method, params)
	//	log4.Debug("Res:%s", data)
	//}
	if err != nil {
		return nil, fmt.Errorf("Call %s error:%s", method, err)
	}
	if err != nil {
		return nil, fmt.Errorf("Call %s error:%s", method, err)
	}
	if data == nil {
		return nil, fmt.Errorf("Call %s return nil.", method)
	}
	res := &DNAJsonRpcRes{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal DNAJsonRpcRes:%s error:%s", res, err)
	}
	data, err = res.HandleResult()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Call sends RPC request to server
func (this *Dna) Call(address string, method string, id interface{}, params []interface{}) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Marshal JSON request: %v\n", err)
		return nil, err
	}
	resp, err := this.client.Post(address, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "POST request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GET response: %v\n", err)
		return nil, err
	}

	return body, nil
}
