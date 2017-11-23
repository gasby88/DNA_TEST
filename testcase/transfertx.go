package testcase

import (
	"DNA/common"
	"DNA/core/transaction"
	"DNA_TEST/dna"
	. "DNA_TEST/testframework"
	"fmt"
	"time"
	"DNA/account"
)

func TestTransferTransaction(ctx *TestFrameworkContext) bool {
	ctx.Dna.WaitForGenerateBlock(10 * time.Second)

	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	//asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}

	programHashTo, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHashTo,
		}
		//dibs output
		output2 := &transaction.TxOutput{
			AssetID:     output.AssetID,
			Value:       unspent.Value - output.Value,
			ProgramHash: unspent.ProgramHash,
		}
		txOutputs = append(txOutputs, output)
		txOutputs = append(txOutputs, output2)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	txHash, err := ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}

	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 30)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return false
	}

	transferTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction TxHash:%x error:%s", txHash, err)
		return false
	}

	txInputs2 := transferTx2.UTXOInputs
	txOutputs2 := transferTx2.Outputs
	ok, err := checkTransferTxResult(txInputs, txInputs2, txOutputs, txOutputs2)
	if err != nil {
		ctx.LogError("checkTransferTxResult error:%s", err)
		return false
	}
	return ok
}

func TestTransferMultiTransaction(ctx *TestFrameworkContext) bool {
	empty := common.Uint256{}
	assetName1 := "TS01"
	assetId1 := ctx.DnaAsset.GetAssetId(assetName1)
	if assetId1 == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName1)
		ctx.FailNow()
		return false
	}
	//asset1 := ctx.DnaAsset.GetAssetByName(assetName1)
	programHash1, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	assetName2 := "TS02"
	assetId2 := ctx.DnaAsset.GetAssetId(assetName2)
	if assetId2 == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName2)
		ctx.FailNow()
		return false
	}
	//asset2 := ctx.DnaAsset.GetAssetByName(assetName2)
	programHash2, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents1, err := ctx.Dna.GetUnspendOutput(assetId1, programHash1)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents1 == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}
	unspents2, err := ctx.Dna.GetUnspendOutput(assetId2, programHash2)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents2 == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}
	txInputs := make([]*transaction.UTXOTxInput, 0)
	txOutputs := make([]*transaction.TxOutput, 0)
	for _, unspent := range unspents1 {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHash2,
		}
		output2 := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       unspent.Value - output.Value,
			ProgramHash: unspent.ProgramHash,
		}
		txOutputs = append(txOutputs, output)
		txOutputs = append(txOutputs, output2)
		break
	}
	for _, unspent := range unspents2 {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHash1,
		}
		output2 := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       unspent.Value - output.Value,
			ProgramHash: programHash2,
		}
		txOutputs = append(txOutputs, output)
		txOutputs = append(txOutputs, output2)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	err = ctx.Dna.SignTransaction(transferTx, []*account.Account{ctx.DnaClient.Account1,ctx.DnaClient.Account2})
	if err != nil {
		ctx.LogError("SignTransaction error:%s", err)
		return  false
	}
	txHash, err := ctx.Dna.DoSendTransaction(transferTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}

	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 30)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return false
	}

	transferTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction TxHash:%x error:%s", txHash, err)
		return false
	}

	txInputs2 := transferTx2.UTXOInputs
	txOutputs2 := transferTx2.Outputs
	ok, err := checkTransferTxResult(txInputs, txInputs2, txOutputs, txOutputs2)
	if err != nil {
		ctx.LogError("checkTransferTxResult error:%s", err)
		return false
	}
	return ok
}

func TestMultiSigTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}
	m := 2
	multiSingers := []*account.Account{ ctx.DnaClient.Account3, ctx.DnaClient.Account4, ctx.DnaClient.Account5}
	programHashTo, err := ctx.Dna.GetAccountsProgramHash(ctx.DnaClient.Account2, m, multiSingers)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHashTo,
		}
		//dibs output
		output2 := &transaction.TxOutput{
			AssetID:     output.AssetID,
			Value:       unspent.Value - output.Value,
			ProgramHash: unspent.ProgramHash,
		}
		txOutputs = append(txOutputs, output)
		txOutputs = append(txOutputs, output2)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	txHash, err := ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}
	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 30)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return false
	}

	transferTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction TxHash:%x error:%s", txHash, err)
		return false
	}

	txInputs2 := transferTx2.UTXOInputs
	txOutputs2 := transferTx2.Outputs
	ok, err := checkTransferTxResult(txInputs, txInputs2, txOutputs, txOutputs2)
	if err != nil {
		ctx.LogError("checkTransferTxResult error:%s", err)
		return false
	}
	if !ok {
		ctx.LogError("checkTransferTxResult failuer")
		return false
	}

	programHashTo, err = ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	txInputs3 := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs3 := make([]*transaction.TxOutput, 0, 1)
	txInput := &transaction.UTXOTxInput{
		ReferTxID:          txHash,
		ReferTxOutputIndex: 0,
	}
	txInputs3 = append(txInputs3, txInput)
	txOutput := &transaction.TxOutput{
		AssetID:     assetId,
		Value:       ctx.Dna.MakeAssetAmount(1),
		ProgramHash: programHashTo,
	}
	txOutputs3 = append(txOutputs3, txOutput)
	transferTx3, err := ctx.Dna.NewTransferAssetTransaction(txInputs3, txOutputs3)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}
	//should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account2, transferTx3)
	if err == nil {
		ctx.LogError("SendTransaction should failed. Sig error")
		return false
	}
	//should failed
	_, err = ctx.Dna.SendMultiSigTransction(ctx.DnaClient.Account2, m, []*account.Account{ctx.DnaClient.Account3}, transferTx3)
	if err == nil {
		ctx.LogError("SendMutiSigTransction should failed. Sig not enough")
		return false
	}
	//should failed
	_, err = ctx.Dna.SendMultiSigTransction(ctx.DnaClient.Account2, m, []*account.Account{ctx.DnaClient.Account1, ctx.DnaClient.Account3}, transferTx3)
	if err == nil {
		ctx.LogError("SendMutiSigTransction should failed. Sig not enough")
		return false
	}
	//should success
	txHash2, err := ctx.Dna.SendMultiSigTransction(ctx.DnaClient.Account2, m, multiSingers, transferTx3)
	if err != nil {
		ctx.LogError("SendMutiSigTransction error:%s", err)
		return false
	}

	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 30)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return false
	}

	transferTx4, err := ctx.Dna.GetTransaction(txHash2)
	if err != nil {
		ctx.LogError("GetMultiSigTransaction TxHash:%x error:%s", txHash2, err)
		return false
	}

	ok, err = checkTransferTxResult(transferTx3.UTXOInputs, transferTx4.UTXOInputs, transferTx3.Outputs, transferTx4.Outputs)
	if err != nil {
		ctx.LogError("checkTransferTxResult error:%s", err)
		return false
	}
	return ok
}

func TestTransferOverAmountTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}

	programHashTo, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       unspent.Value * 2,
			ProgramHash: programHashTo,
		}
		txOutputs = append(txOutputs, output)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err == nil || err.Error() != dna.DnaRpcInternalError {
		ctx.LogError("SendTransaction should failed.")
		return false
	}
	return true
}

func TestTransferNegAmountTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	//asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}

	programHashTo, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(-1),
			ProgramHash: programHashTo,
		}
		txOutputs = append(txOutputs, output)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err == nil || err.Error() != dna.DnaRpcInternalError {
		ctx.LogError("SendTransaction should failed")
		return false
	}
	return true
}

func TestTransferPreciseTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	//asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}

	programHashTo, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1.00001),
			ProgramHash: programHashTo,
		}
		txOutputs = append(txOutputs, output)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err == nil || err.Error() != dna.DnaRpcInternalError {
		ctx.LogError("SendTransaction should failed")
		return false
	}
	return true
}

func TestTransferDoubleSpendTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	//asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}

	programHashTo, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0)
	txOutputs := make([]*transaction.TxOutput, 0)
	var txUnspent *dna.UnspendUTXO
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		txUnspent = unspent
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHashTo,
		}
		output2 := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       unspent.Value - output.Value,
			ProgramHash: unspent.ProgramHash,
		}
		txOutputs = append(txOutputs, output)
		txOutputs = append(txOutputs, output2)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	ctx.LogInfo("TestTransferDoubleSpendTransaction part 2")
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}

	//ctx.LogInfo("WaitForGenerateBlock")
	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 30)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false
	}

	dupInputs := txInputs
	dupOutput := &transaction.TxOutput{
		AssetID:     assetId,
		Value:       txUnspent.Value,
		ProgramHash: programHashTo,
	}
	dupOutputs := []*transaction.TxOutput{dupOutput}
	transferTx2, err := ctx.Dna.NewTransferAssetTransaction(dupInputs, dupOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx2)
	if err == nil || err.Error() != dna.DnaRpcInternalError {
		ctx.LogError("SendTransaction should failed")
		return false
	}

	ctx.LogInfo("Test mutil utxo double spend")
	unspents, err = ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}
	txInputs2 := make([]*transaction.UTXOTxInput, 0)
	txOutputs2 := make([]*transaction.TxOutput, 0)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs2 = append(txInputs2, input)
		txInputs2 = append(txInputs2, txInputs[0]) //duplicate
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHashTo,
		}
		output2 := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       unspent.Value - output.Value,
			ProgramHash: unspent.ProgramHash,
		}
		txOutputs2 = append(txOutputs2, output)
		txOutputs2 = append(txOutputs2, output2)
		txOutputs2 = append(txOutputs2, dupOutput)
		break
	}

	transferTx, err = ctx.Dna.NewTransferAssetTransaction(txInputs2, txOutputs2)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err == nil || err.Error() != dna.DnaRpcInternalError {
		ctx.LogError("SendTransaction should failed")
		return false
	}

	return true
}

func TestTransferDuplicateUTXOTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	//asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}

	programHashTo, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		txInputs = append(txInputs, input) //Duplicate UTXO Input
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHashTo,
		}
		txOutputs = append(txOutputs, output)
		output2 := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       unspent.Value*2 - output.Value,
			ProgramHash: unspent.ProgramHash,
		}
		txOutputs = append(txOutputs, output2)
		break
	}
	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}
	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, transferTx)
	if err == nil {
		ctx.LogError("SendTransaction should failed")
		return false
	}
	return true
}

func TestTransferInvalidAccountTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetName)
		ctx.FailNow()
		return false
	}
	//asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}
	unspents, err := ctx.Dna.GetUnspendOutput(assetId, programHash)
	if err != nil {
		ctx.LogError("GetUnspendOutput error:%s", err)
		return false
	}
	if unspents == nil {
		ctx.LogError("GetUnspendOutput return nil")
		return false
	}

	programHashTo, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account2)
	if err != nil {
		ctx.LogError("GetAccountProgramHash error:%s", err)
		return false
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	for _, unspent := range unspents {
		if unspent.Value < 1 {
			continue
		}
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)
		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       ctx.Dna.MakeAssetAmount(1),
			ProgramHash: programHashTo,
		}
		txOutputs = append(txOutputs, output)
		break
	}
	if len(txInputs) == 0 {
		ctx.LogError("TxInput is nil")
		return false
	}

	transferTx, err := ctx.Dna.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account2, transferTx)
	if err == nil || err.Error() != dna.DnaRpcInternalError {
		ctx.LogError("SendTransaction should failed")
		return false
	}
	return true
}

func checkTransferTxResult(txInputs, txInputsRes []*transaction.UTXOTxInput,
	txOutputs, txOutputsRes []*transaction.TxOutput) (bool, error) {
	if len(txInputsRes) != len(txInputs) ||
		len(txOutputsRes) != len(txOutputs) {
		return false, fmt.Errorf("len(txInputs2):%v != len(txInputs):%v or len(txOutputs2):%v != len(txOutputs):%v ",
			len(txInputsRes),
			len(txInputs),
			len(txOutputsRes),
			len(txOutputs))
	}

	for i, txInputRes := range txInputsRes {
		txInput := txInputs[i]
		if txInput.ReferTxOutputIndex != txInputRes.ReferTxOutputIndex ||
			txInput.ReferTxID != txInputRes.ReferTxID {
			return false, fmt.Errorf("ReferTxID:%x != %x or ReferTxOutputIndex:%v != %v",
				txInputRes.ReferTxID,
				txInput.ReferTxID,
				txInputRes.ReferTxOutputIndex,
				txInput.ReferTxOutputIndex)
		}
	}

	for i, txOutputRes := range txOutputsRes {
		txOutput := txOutputs[i]
		if txOutput.ProgramHash != txOutputRes.ProgramHash ||
			txOutput.Value != txOutputRes.Value ||
			txOutput.AssetID != txOutputRes.AssetID {
			return false, fmt.Errorf("ProgramHash:%x != %x or Value:%v != %v or AssetID:%x != %x",
				txOutputRes.ProgramHash,
				txOutput.ProgramHash,
				txOutputRes.Value,
				txOutput.Value,
				txOutputRes.ProgramHash,
				txOutput.ProgramHash)
		}
	}
	return true, nil
}
