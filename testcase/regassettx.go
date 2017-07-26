package testcase

import (
	"DNA/common"
	. "DNA/core/asset"
	"DNA/core/transaction/payload"
	. "DNA_TEST/dna"
	. "DNA_TEST/testframework"
	"time"
)

func TestRegisterAssetTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetPrecise := byte(4)
	assetType := Token
	recordType := UTXO
	asset := ctx.Dna.CreateAsset(assetName, assetPrecise, assetType, recordType)
	assetAmount := ctx.Dna.MakeAssetAmount(20000)
	if !testRegisterAssetTransaction(asset, assetAmount, ctx) {
		ctx.LogError("TestRegisterAssetTransaction Asset:%+v Amount:%v test failed.",
			asset, assetAmount)
		ctx.FailNow()
		return false
	}

	assetName = "TS02"
	assetPrecise = byte(8)
	assetType = Share
	recordType = UTXO
	asset = ctx.Dna.CreateAsset(assetName, assetPrecise, assetType, recordType)
	assetAmount = ctx.Dna.MakeAssetAmount(100000)
	if !testRegisterAssetTransaction(asset, assetAmount, ctx) {
		ctx.LogError("TestRegisterAssetTransaction Asset:%+v Amount:%v test failed.",
		asset, assetAmount)
		ctx.FailNow()
		return false
	}

	return true
}

func TestRegisterAssetPreciseTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetPrecise := byte(4)
	assetType := Token
	recordType := UTXO
	asset := ctx.Dna.CreateAsset(assetName, assetPrecise, assetType, recordType)
	assetAmount := ctx.Dna.MakeAssetAmount(100.00001)
	regTx, err := ctx.Dna.NewAssetRegisterTransaction(asset, assetAmount, ctx.DnaClient.Admin, ctx.DnaClient.Admin)
	if err != nil {
		ctx.LogError("NewAssetRegisterTransaction Asset:%+v Amount:%v Admin:%+v Account:%+v error:%s",
			asset,
			assetAmount,
			ctx.DnaClient.Admin,
			ctx.DnaClient.Admin,
			err)

		ctx.FailNow()
		return false
	}

	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Admin, regTx)
	if err == nil || err.Error() != DnaRpcInternalError {
		ctx.LogError("SendTransaction AssetRegisterTransaction should failed error:%s", err)
		return false
	}
	return true
}

func TestRegisterAssetMaxPreciseTransaction(ctx *TestFrameworkContext) bool {
	assetName := "TS01"
	assetPrecise := byte(9)
	assetType := Token
	recordType := UTXO
	asset := ctx.Dna.CreateAsset(assetName, assetPrecise, assetType, recordType)
	assetAmount := ctx.Dna.MakeAssetAmount(10000)

	regTx, err := ctx.Dna.NewAssetRegisterTransaction(asset, assetAmount, ctx.DnaClient.Admin, ctx.DnaClient.Admin)
	if err != nil {
		ctx.LogError("NewAssetRegisterTransaction Asset:%+v Amount:%v Admin:%+v Account:%+v error:%s",
			asset,
			assetAmount,
			ctx.DnaClient.Admin,
			ctx.DnaClient.Admin,
			err)
		return false
	}
	//Should failed
	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Admin, regTx)
	if err == nil || err.Error() != DnaRpcInternalError {
		ctx.LogError("SendTransaction AssetRegisterTransaction should failed err:s", err)
		return false
	}
	return true
}

func testRegisterAssetTransaction(asset *Asset, assetAmount common.Fixed64, ctx *TestFrameworkContext) bool {
	regTx, err := ctx.Dna.NewAssetRegisterTransaction(asset, assetAmount, ctx.DnaClient.Admin, ctx.DnaClient.Admin)
	if err != nil {
		ctx.LogError("NewAssetRegisterTransaction Asset:%+v Amount:%v Admin:%+v Account:%+v error:%s",
			asset,
			assetAmount,
			ctx.DnaClient.Admin,
			ctx.DnaClient.Admin,
			err)

		ctx.FailNow()
		return false
	}

	txHash, err := ctx.Dna.SendTransaction(ctx.DnaClient.Admin, regTx)
	if err != nil {
		ctx.LogError("SendTransaction AssetRegisterTransaction error:%s", err)
		ctx.FailNow()
		return false
	}

	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 30)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return false
	}

	regTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction Hash:%x error:%s", txHash, err)
		return false
	}

	regAssetPayload := regTx2.Payload.(*payload.RegisterAsset)
	asset2 := regAssetPayload.Asset
	if !AssetEqualTo(asset, asset2) || regAssetPayload.Amount != assetAmount {
		ctx.LogError("Asset get from transaction not equal.")
		return false
	}

	ctx.DnaAsset.RegAsset(txHash, asset)
	return true
}
