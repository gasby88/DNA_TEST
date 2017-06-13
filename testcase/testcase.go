package testcase

import(
	. "DNA_TEST/testframework"
)

func init(){
	TFramework.RegTestCase("TestRegisterAssetTransaction", TestRegisterAssetTransaction)
	TFramework.RegTestCase("TestRecordTransactionByRecord", TestRecordTransactionByRecord)
	TFramework.RegTestCase("TestRecordTransactionByTransfer", TestRecordTransactionByTransfer)
	TFramework.RegTestCase("TestIssueAssetTransaction", TestIssueAssetTransaction)
	TFramework.RegTestCase("TestIssueAssetMutiTransaction", TestIssueAssetMutiTransaction)
	TFramework.RegTestCase("TestTransferTransaction",TestTransferTransaction)
	TFramework.RegTestCase("TestTransferMutiTransaction",TestTransferMutiTransaction)
}
