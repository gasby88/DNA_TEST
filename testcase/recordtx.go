package testcase

import (
	"DNA/core/transaction"
	. "DNA/core/transaction/payload"
	. "DNA_TEST/testframework"
	"time"
)

func TestRecordTransactionByRecord(ctx *TestFrameworkContext) bool {
	recordType := "TestRecord"
	recordData := []byte("Hello World!")
	recordTx, err := ctx.Dna.NewRecordTransaction(recordType, recordData)
	if err != nil {
		ctx.LogError("NewRecordTransaction RecordType:%s RecordData:%s error:%s", recordType, recordData, err)
		return false
	}

	txHash, err := ctx.Dna.SendTransaction(ctx.DnaClient.Account1, recordTx)
	if err != nil {
		ctx.LogError("SendTransaction RecordTransaction error:%s", err)
		return false
	}

	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 10)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false
	}

	recordTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction error:%s", err)
		return false
	}

	testRecord := recordTx2.Payload.(*Record)
	if testRecord.RecordType != recordType || string(testRecord.RecordData) != string(recordData) {
		ctx.LogError("RecordType:%s != %s or RecordData:%s != %s", testRecord.RecordType, recordType, testRecord.RecordData, recordData)
		return false
	}

	return true
}

func TestRecordTransactionByTransfer(ctx *TestFrameworkContext) bool {
	recordData := []byte("Hello World!")

	recordTx, err := ctx.Dna.NewTransferAssetTransaction(nil, nil)
	if err != nil {
		ctx.LogError("NewTransferAssetTransaction error:%s", err)
		return false
	}

	attrType := transaction.Description
	attrs := &transaction.TxAttribute{
		Usage: transaction.Description,
		Data:  recordData,
		Size:  uint32(len(recordData)),
	}
	recordTx.Attributes = append(recordTx.Attributes, attrs)

	txHash, err := ctx.Dna.SendTransaction(ctx.DnaClient.Account1, recordTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}
	ctx.LogInfo("TxHash:%x", txHash)
	_, err = ctx.Dna.WaitForGenerateBlock(time.Second * 10)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false
	}

	recordTx2, err := ctx.Dna.GetTransaction(txHash)
	if err != nil {
		ctx.LogError("GetTransaction error:%s", err)
		return false
	}

	recordData2 := []byte("")
	for _, attr := range recordTx2.Attributes{
		if attr.Usage == attrType{
			recordData2 = attr.Data
		}
	}

	if string(recordData) != string(recordData2) {
		ctx.LogError("RecordData:%s != %s", recordData2, recordData)
		return false
	}

	return true
}
