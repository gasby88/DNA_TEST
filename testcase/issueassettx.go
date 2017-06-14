package testcase

import (
	"DNA/common"
	"DNA/core/transaction"
	"DNA_TEST/testframework"
	."DNA_TEST/dna"
	"fmt"
	"time"
)

func TestIssueAssetTransaction(ctx *testframework.TestFrameworkContext) bool {
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist",assetName)
		ctx.FailNow()
		return false
	}
	asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetProgramHash error:%s", err)
		return false
	}
	output := &transaction.TxOutput{
		Value:       ctx.Dna.MakeAssetAmount(100, asset.Precision),
		AssetID:     assetId,
		ProgramHash: programHash,
	}
	txOutputs := []*transaction.TxOutput{output}
	issueTx, err := ctx.Dna.NewIssueAssetTransaction(txOutputs)
	if err != nil {
		ctx.LogError("NewIssueAssetTransaction error:%s", err)
		return false
	}
	txHash, err := ctx.Dna.SendTransaction(ctx.DnaClient.Admin, issueTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}

	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 10)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return false
	}
	issueTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction TxHash:%x error:%s", txHash, err)
		return false
	}
	if len(issueTx2.Outputs) == 0 {
		ctx.LogError("GetTransaction Outputs error")
		return false
	}

	txOutputsRes := issueTx2.Outputs
	ok, err := checkIssueAssetTxResult(txOutputs, txOutputsRes)
	if err != nil {
		ctx.LogError("checkIssueAssetTxResult error:%s", err)
		return false
	}
	return ok
}

func TestIssueAssetMutiTransaction(ctx *testframework.TestFrameworkContext) bool {
	empty := common.Uint256{}
	assetName1 := "TS01"
	assetId1 := ctx.DnaAsset.GetAssetId(assetName1)
	if assetId1 == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetId1)
		ctx.FailNow()
		return false
	}
	asset1 := ctx.DnaAsset.GetAssetByName(assetName1)
	assetName2 := "TS02"
	assetId2 := ctx.DnaAsset.GetAssetId(assetName2)
	if assetId1 == empty {
		ctx.LogError("AssetName:%s doesnot exist", assetId2)
		ctx.FailNow()
		return false
	}
	asset2 := ctx.DnaAsset.GetAssetByName(assetName2)
	programHash1, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetProgramHash error:%s", err)
		return false
	}
	programHash2, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetProgramHash error:%s", err)
		return false
	}
	txOutput1 := &transaction.TxOutput{
		Value:       ctx.Dna.MakeAssetAmount(100, asset1.Precision),
		AssetID:     assetId1,
		ProgramHash: programHash1,
	}
	txOutput2 := &transaction.TxOutput{
		Value:       ctx.Dna.MakeAssetAmount(100, asset2.Precision),
		AssetID:     assetId2,
		ProgramHash: programHash2,
	}
	txOutput3 := &transaction.TxOutput{
		Value:       ctx.Dna.MakeAssetAmount(100, asset2.Precision),
		AssetID:     assetId2,
		ProgramHash: programHash1,
	}
	txOutputs := []*transaction.TxOutput{txOutput1, txOutput2, txOutput3}
	issueTx, err := ctx.Dna.NewIssueAssetTransaction(txOutputs)
	if err != nil {
		ctx.LogError("NewIssueAssetTransaction error:%s", err)
		return false
	}
	txHash, err := ctx.Dna.SendTransaction(ctx.DnaClient.Admin, issueTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}

	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 10)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return false
	}
	issueTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction TxHash:%x error:%s", txHash, err)
		return false
	}
	if len(issueTx2.Outputs) == 0 {
		ctx.LogError("GetTransaction Outputs error")
		return false
	}

	txOutputsRes := issueTx2.Outputs
	ok, err := checkIssueAssetTxResult(txOutputs, txOutputsRes)
	if err != nil {
		ctx.LogError("checkIssueAssetTxResult error:%s", err)
		return false
	}
	return ok
}

func TestIssueAssetOverAmountTransaction(ctx *testframework.TestFrameworkContext) bool{
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist",assetName)
		ctx.FailNow()
		return false
	}
	asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetProgramHash error:%s", err)
		return false
	}
	output := &transaction.TxOutput{
		Value:       ctx.Dna.MakeAssetAmount(1000000000000000, asset.Precision),
		AssetID:     assetId,
		ProgramHash: programHash,
	}
	txOutputs := []*transaction.TxOutput{output}
	issueTx, err := ctx.Dna.NewIssueAssetTransaction(txOutputs)
	if err != nil {
		ctx.LogError("NewIssueAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Admin, issueTx)
	if err == nil || err.Error() != DnaRpcInternalError  {
		ctx.LogError("SendTransaction should failed.")
		return false
	}
	return true
}

func TestIssueAssetNegAmountTransaction(ctx *testframework.TestFrameworkContext)bool{
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist",assetName)
		ctx.FailNow()
		return false
	}
	asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetProgramHash error:%s", err)
		return false
	}
	output := &transaction.TxOutput{
		Value:       ctx.Dna.MakeAssetAmount(-1, asset.Precision),
		AssetID:     assetId,
		ProgramHash: programHash,
	}
	txOutputs := []*transaction.TxOutput{output}
	issueTx, err := ctx.Dna.NewIssueAssetTransaction(txOutputs)
	if err != nil {
		ctx.LogError("NewIssueAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Admin, issueTx)
	if err == nil || err.Error() != DnaRpcInternalError  {
		ctx.LogError("SendTransaction error should failed.")
		return false
	}
	return true
}

func TestIssueAssetPreciseTransaction(ctx *testframework.TestFrameworkContext)bool{
	assetName := "TS01"
	assetId := ctx.DnaAsset.GetAssetId(assetName)
	empty := common.Uint256{}
	if assetId == empty {
		ctx.LogError("AssetName:%s doesnot exist",assetName)
		ctx.FailNow()
		return false
	}
	asset := ctx.DnaAsset.GetAssetByName(assetName)
	programHash, err := ctx.Dna.GetAccountProgramHash(ctx.DnaClient.Account1)
	if err != nil {
		ctx.LogError("GetProgramHash error:%s", err)
		return false
	}
	output := &transaction.TxOutput{
		Value:       ctx.Dna.MakeAssetAmount(100.00001, asset.Precision),
		AssetID:     assetId,
		ProgramHash: programHash,
	}
	txOutputs := []*transaction.TxOutput{output}
	issueTx, err := ctx.Dna.NewIssueAssetTransaction(txOutputs)
	if err != nil {
		ctx.LogError("NewIssueAssetTransaction error:%s", err)
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Admin, issueTx)
	if err == nil || err.Error() != DnaRpcInternalError  {
		ctx.LogError("SendTransaction error.Transaction shuld be rejected")
		return false
	}
	return true
}

func checkIssueAssetTxResult(txOutputs, txOutputsRes []*transaction.TxOutput) (bool, error) {
	if len(txOutputs) != len(txOutputsRes) {
		return false, fmt.Errorf("len(txOutputs):%v != len(txOutputsRes):%v", len(txOutputs), len(txOutputsRes))
	}
	for i, txOutputRes := range txOutputsRes {
		txOutput := txOutputs[i]
		if txOutput.ProgramHash != txOutputRes.ProgramHash &&
			txOutput.AssetID != txOutputRes.AssetID &&
			txOutput.Value != txOutputRes.Value {
			return false, fmt.Errorf("IssueAssetTransaction ProgramHash:%x != %x AssetID:%x != %x Value:%v != %v",
				txOutputRes.ProgramHash,
				txOutput.ProgramHash,
				txOutputRes.AssetID,
				txOutput.AssetID,
				txOutputRes.Value,
				txOutput.Value,
			)
		}
	}
	return true, nil
}
