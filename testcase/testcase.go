package testcase

import (
	. "DNA_TEST/testframework"
)

func init() {
	//Register Asset test
	TFramework.RegTestCase("TestRegisterAssetTransaction", TestRegisterAssetTransaction)
	TFramework.RegTestCase("TestRegisterAssetPreciseTransaction", TestRegisterAssetPreciseTransaction)
	TFramework.RegTestCase("TestRegisterAssetMaxPreciseTransaction", TestRegisterAssetMaxPreciseTransaction)

	//Issue Asset test
	TFramework.RegTestCase("TestIssueAssetTransaction", TestIssueAssetTransaction)
	TFramework.RegTestCase("TestIssueAssetMutiTransaction", TestIssueAssetMutiTransaction)
	TFramework.RegTestCase("TestIssueAssetNegAmountTransaction", TestIssueAssetNegAmountTransaction)
	TFramework.RegTestCase("TestIssueAssetPreciseTransaction", TestIssueAssetPreciseTransaction)
	TFramework.RegTestCase("TestIssueAssetOverAmountTransaction", TestIssueAssetOverAmountTransaction)

	//Transfer Asset test
	TFramework.RegTestCase("TestTransferTransaction", TestTransferTransaction)
	TFramework.RegTestCase("TestTransferMutiTransaction", TestTransferMutiTransaction)
	TFramework.RegTestCase("TestTransferOverAmountTransaction", TestTransferOverAmountTransaction)
	TFramework.RegTestCase("TestTransferNegAmountTransaction", TestTransferNegAmountTransaction)
	TFramework.RegTestCase("TestTransferPreciseTransaction", TestTransferPreciseTransaction)
	TFramework.RegTestCase("TestTransferDoubleSendTransaction", TestTransferDoubleSpendTransaction)
	TFramework.RegTestCase("TestTransferInvalidAccountTransaction", TestTransferInvalidAccountTransaction)
	TFramework.RegTestCase("TestTransferDuplicateUTXOTransaction", TestTransferDuplicateUTXOTransaction)

	//Record test
	TFramework.RegTestCase("TestRecordTransactionByRecord", TestRecordTransactionByRecord)
	TFramework.RegTestCase("TestRecordTransactionByTransfer", TestRecordTransactionByTransfer)

	//Benchmark
	TFramework.RegBenchTestCase("BenchmarkTransaction", BenchmarkTransaction)
}
