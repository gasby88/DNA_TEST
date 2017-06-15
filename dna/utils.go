package dna

import (
	"DNA/common"
	"DNA/core/asset"
	"DNA/core/code"
	"DNA/core/contract"
	"DNA/core/contract/program"
	"DNA/core/ledger"
	"DNA/core/transaction"
	txpl "DNA/core/transaction/payload"
	"DNA/crypto"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	//"os"
)

func ParseTransaction(txStr *Transactions) (*transaction.Transaction, error) {
	payload, err := ParseToPayload(txStr.TxType, []byte(txStr.Payload))
	if err != nil {
		return nil, fmt.Errorf("ParseToPayload:%s error:%s", txStr.Payload, err)
	}

	attris := make([]*transaction.TxAttribute, len(txStr.Attributes))
	for i, attr := range txStr.Attributes {
		txAttr, err := ParseTransactionAttributes(&attr)
		if err != nil {
			return nil, fmt.Errorf("ParseTransactionAttributes:%+v error:%s", attr, err)
		}
		attris[i] = txAttr
	}

	utxoInputs := make([]*transaction.UTXOTxInput, len(txStr.UTXOInputs))
	for i, input := range txStr.UTXOInputs {
		txInput, err := ParseTransactionUTXOTxInput(&input)
		if err != nil {
			return nil, fmt.Errorf("ParseTransactionUTXOTxInput:%+v error:%s", input, err)
		}
		utxoInputs[i] = txInput
	}

	balance := make([]*transaction.BalanceTxInput, len(txStr.BalanceInputs))
	for i, input := range txStr.BalanceInputs {
		txInput, err := ParseTransactionBalanceTxInputInfo(&input)
		if err != nil {
			return nil, fmt.Errorf("ParseTransactionBalanceTxInputInfo:%+v error:%s", input, err)
		}
		balance[i] = txInput
	}

	outputs := make([]*transaction.TxOutput, len(txStr.Outputs))
	for i, output := range txStr.Outputs {
		txOutput, err := ParseTransactionOutputs(&output)
		if err != nil {
			return nil, fmt.Errorf("ParseTransactionOutputs:%+v error:%s", output, err)
		}
		outputs[i] = txOutput
	}

	programs := make([]*program.Program, len(txStr.Programs))
	for i, p := range txStr.Programs {
		txProgram, err := ParseTransactionPrograms(&p)
		if err != nil {
			return nil, fmt.Errorf("ParseTransactionPrograms:%+v error:%s", p, err)
		}
		programs[i] = txProgram
	}

	assetOutputs := make(map[common.Uint256][]*transaction.TxOutput, len(txStr.AssetOutputs))
	for _, assetOutput := range txStr.AssetOutputs {
		outputs := make([]*transaction.TxOutput, len(assetOutput.Txout))
		for i, output := range assetOutput.Txout {
			txOutput, err := ParseTransactionOutputs(&output)
			if err != nil {
				return nil, fmt.Errorf("AssetOutputs ParseTransactionOutputs:%+v error:%s", output, err)
			}
			outputs[i] = txOutput
		}
		assetOutputs[assetOutput.Key] = outputs
	}

	assetInputAmounts := make(map[common.Uint256]common.Fixed64, len(txStr.AssetInputAmount))
	for _, assetInputAmount := range txStr.AssetInputAmount {
		assetInputAmounts[assetInputAmount.Key] = assetInputAmount.Value
	}

	assetOutputAmounts := make(map[common.Uint256]common.Fixed64, len(txStr.AssetOutputAmount))
	for _, assetOutputAmount := range txStr.AssetOutputAmount {
		assetOutputAmounts[assetOutputAmount.Key] = assetOutputAmount.Value
	}

	tx := &transaction.Transaction{}
	tx.TxType = transaction.TransactionType(txStr.TxType)
	tx.PayloadVersion = txStr.PayloadVersion
	tx.Nonce = txStr.Nonce
	tx.Payload = payload
	tx.AssetOutputAmount = assetOutputAmounts
	tx.AssetInputAmount = assetInputAmounts
	tx.AssetOutputs = assetOutputs
	tx.Programs = programs
	tx.Outputs = outputs
	tx.BalanceInputs = balance
	tx.Attributes = attris
	tx.UTXOInputs = utxoInputs

	txHash, err := ParseUint256FromString(txStr.Hash)
	if err != nil {
		return nil, fmt.Errorf("Hash ParseUint256FromString:%s error:%s", txStr.Hash, err)
	}
	tx.SetHash(txHash)
	return tx, nil
}

func ParseToPayload(payloadType transaction.TransactionType, data json.RawMessage) (transaction.Payload, error) {
	var payload transaction.Payload

	switch payloadType {
	case transaction.RegisterAsset:
		p := &PayloadRegisterAssetInfo{}
		err := json.Unmarshal(data, p)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal payload RegisterAssetInfo:%s error:%s", data, err)
		}
		regAsset, err := ParseRegisterAssetInfo(p)
		if err != nil {
			return nil, fmt.Errorf("ParsePayloadRegisterAssetInfo error:%s", err)
		}
		payload = regAsset
	case transaction.Record:
		p := &PayloadRecord{}
		err := json.Unmarshal(data, p)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal payload Record:%s error:%s", data, err)
		}
		record, err := ParseRecord(p)
		if err != nil {
			return nil, fmt.Errorf("ParsePayloadRecord error:%s", err)
		}
		payload = record
	case transaction.DeployCode:
		p := &PayloadDeployCodeInfo{}
		err := json.Unmarshal(data, p)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal payload DeployCodeInfo:%s error:%s", data, err)
		}

		deplyCode, err := ParseDeployCodeInfo(p)
		if err != nil {
			return nil, fmt.Errorf("ParsePayloadDeployCodeInfo error:%s", err)
		}
		payload = deplyCode
	}

	return payload, nil
}

func ParseTransactionAttributes(attr *TxAttributeInfo) (*transaction.TxAttribute, error) {
	data, err := hex.DecodeString(attr.Date)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString TxAttributeInfo.Data:%s error:%s", attr.Date, err)
	}
	txAttr := &transaction.TxAttribute{}
	txAttr.Size = attr.Size
	txAttr.Usage = transaction.TransactionAttributeUsage(attr.Usage)
	txAttr.Data = data
	return txAttr, nil
}

func ParseTransactionUTXOTxInput(input *UTXOTxInputInfo) (*transaction.UTXOTxInput, error) {
	txId, err := ParseUint256FromString(input.ReferTxID)
	if err != nil {
		return nil, fmt.Errorf("ParseUint256FromString UTXOTxInputInfo.ReferTxID:%s error:%s", input.ReferTxID, err)
	}
	return &transaction.UTXOTxInput{
		ReferTxID:          txId,
		ReferTxOutputIndex: input.ReferTxOutputIndex,
	}, nil
}

func ParseTransactionBalanceTxInputInfo(input *BalanceTxInputInfo) (*transaction.BalanceTxInput, error) {
	assetId, err := ParseUint256FromString(input.AssetID)
	if err != nil {
		return nil, fmt.Errorf("ParseUint256FromString BalanceTxInputInfo.AssetID:%s error:%s", input.AssetID, err)
	}
	programHash, err := ParseUint160FromString(input.ProgramHash)
	if err != nil {
		return nil, fmt.Errorf("ParseUint160FromString BalanceTxInputInfo.ProgramHash:%s error:%s", input.ProgramHash, err)
	}
	return &transaction.BalanceTxInput{
		AssetID:     assetId,
		Value:       input.Value,
		ProgramHash: programHash,
	}, nil
}

func ParseTransactionOutputs(output *TxoutputInfo) (*transaction.TxOutput, error) {
	assetId, err := ParseUint256FromString(output.AssetID)
	if err != nil {
		return nil, fmt.Errorf("ParseUint256FromString TxOutput.AssetID:%s error:%s", output.AssetID, err)
	}
	programHash, err := ParseUint160FromString(output.ProgramHash)
	if err != nil {
		return nil, fmt.Errorf("ParseUint160FromString TxOutput.ProgramHash:%s error:%s", output.ProgramHash, err)
	}
	return &transaction.TxOutput{
		AssetID:     assetId,
		Value:       output.Value,
		ProgramHash: programHash,
	}, nil
}

func ParseTransactionPrograms(p *ProgramInfo) (*program.Program, error) {
	code, err := hex.DecodeString(p.Code)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString Code:%s error:%s", p.Code, err)
	}
	param, err := hex.DecodeString(p.Parameter)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString Parameter:%s error:%s", p.Parameter, err)
	}
	return &program.Program{
		Code:      code,
		Parameter: param,
	}, nil
}

func ParseRegisterAssetInfo(p *PayloadRegisterAssetInfo) (*txpl.RegisterAsset, error) {
	regAsset := &txpl.RegisterAsset{}
	regAsset.Asset = p.Asset
	regAsset.Amount = p.Amount

	controler, err := ParseUint160FromString(p.Controller)
	if err != nil {
		return nil, fmt.Errorf("Controller:%s ParseUint160FromString error:%s", p.Controller, err)
	}
	regAsset.Controller = controler

	x := &big.Int{}
	_, err = fmt.Sscan(p.Issuer.X, x)
	if err != nil {
		return nil, fmt.Errorf("fmt.Sscan Issuer.X:%s error:%s", p.Issuer.X, err)
	}
	y := &big.Int{}
	_, err = fmt.Sscan(p.Issuer.Y, y)
	if err != nil {
		return nil, fmt.Errorf("fmt.Sscan Issuer.Y:%s error:%s", p.Issuer.Y, err)
	}

	issuer := &crypto.PubKey{
		X: x,
		Y: y,
	}
	regAsset.Issuer = issuer
	return regAsset, nil
}

func ParseRecord(p *PayloadRecord) (*txpl.Record, error) {
	record := &txpl.Record{}
	record.RecordType = p.RecordType
	data, err := hex.DecodeString(p.RecordData)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString RecordData:%s error:%s", p.RecordData, err)
	}

	record.RecordData = data
	return record, nil
}

func ParseDeployCodeInfo(p *PayloadDeployCodeInfo) (*txpl.DeployCode, error) {
	c, err := hex.DecodeString(p.Code.Code)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString Code:%s error:%s", p.Code.Code, err)
	}
	paramByte, err := hex.DecodeString(p.Code.ParameterTypes)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString ParameterTypes:%s error:%s", p.Code.ParameterTypes, err)
	}
	param := contract.ByteToContractParameterType(paramByte)
	retByte, err := hex.DecodeString(p.Code.ReturnTypes)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString ReturnTypes:%s error:%s", p.Code.ReturnTypes, err)
	}
	ret := contract.ByteToContractParameterType(retByte)

	deplyCode := &txpl.DeployCode{}
	deplyCode.Code = &code.FunctionCode{
		Code:           c,
		ParameterTypes: param,
		ReturnTypes:    ret,
	}
	deplyCode.Name = p.Name
	deplyCode.Author = p.Author
	deplyCode.CodeVersion = p.CodeVersion
	deplyCode.Description = p.Description
	deplyCode.Email = p.Email
	return deplyCode, nil
}

func ParseBlock(blockInfo *BlockInfo) (*ledger.Block, error) {
	txs := make([]*transaction.Transaction, len(blockInfo.Transactions))
	for i, txStr := range blockInfo.Transactions {
		tx, err := ParseTransaction(txStr)
		if err != nil {
			return nil, fmt.Errorf("ParseTransaction transactions:%s error:%s", txStr, err)
		}
		txs[i] = tx
	}

	program, err := ParseTransactionPrograms(&blockInfo.BlockData.Program)
	if err != nil {
		return nil, fmt.Errorf("ParseTransactionPrograms Program:%s error:%s", blockInfo.BlockData.Program, err)
	}
	nextBookKeeper, err := ParseUint160FromString(blockInfo.BlockData.NextBookKeeper)
	if err != nil {
		return nil, fmt.Errorf("ParseUint160FromString NextBookKeeper:%s error:%s", blockInfo.BlockData.NextBookKeeper, err)
	}
	prevBlockHash, err := ParseUint256FromString(blockInfo.BlockData.PrevBlockHash)
	if err != nil {
		return nil, fmt.Errorf("ParseUint256FromString PrevBlockHash:%s error:%s", blockInfo.BlockData.PrevBlockHash, err)
	}
	txRoot, err := ParseUint256FromString(blockInfo.BlockData.TransactionsRoot)
	if err != nil {
		return nil, fmt.Errorf("ParseUint256FromString TransactionsRoot:%s error:%s", blockInfo.BlockData.TransactionsRoot, err)
	}
	blockHead := &ledger.Blockdata{}
	blockHead.Program = program
	blockHead.NextBookKeeper = nextBookKeeper
	blockHead.Height = blockInfo.BlockData.Height
	blockHead.Timestamp = blockInfo.BlockData.Timestamp
	blockHead.Version = blockInfo.BlockData.Version
	blockHead.PrevBlockHash = prevBlockHash
	blockHead.ConsensusData = blockInfo.BlockData.ConsensusData
	blockHead.TransactionsRoot = txRoot

	return &ledger.Block{
		Blockdata:    blockHead,
		Transactions: txs,
	}, nil
}

func ParseUint160FromString(value string) (common.Uint160, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return common.Uint160{}, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	res, err := common.Uint160ParseFromBytes(data)
	if err != nil {
		return common.Uint160{}, fmt.Errorf("Uint160ParseFromBytes error:%s", err)
	}
	return res, nil
}

func ParseUint256FromString(value string) (common.Uint256, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	res, err := common.Uint256ParseFromBytes(data)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("Uint160ParseFromBytes error:%s", err)
	}
	return res, nil
}

func Uint256ToString(value common.Uint256) string {
	return hex.EncodeToString(value.ToArray())
}

func Uint160ToString(value common.Uint160) string {
	return hex.EncodeToString(value.ToArray())
}

//func FileExisted(filename string) bool {
//	_, err := os.Stat(filename)
//	return err == nil || os.IsExist(err)
//}

func AssetEqualTo(as1, as2 *asset.Asset) bool {
	if as1 == nil && as2 == nil {
		return true
	}
	if as1 == nil || as2 == nil {
		return false
	}
	if as1.Name == as1.Name &&
		as1.RecordType == as2.RecordType &&
		as1.AssetType == as2.AssetType &&
		as1.Precision == as2.Precision {
		return true
	}
	return false
}

type ByProgramHashes []common.Uint160

func (a ByProgramHashes) Len() int      { return len(a) }
func (a ByProgramHashes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByProgramHashes) Less(i, j int) bool {
	if a[i].CompareTo(a[j]) > 0 {
		return false
	} else {
		return true
	}
}
