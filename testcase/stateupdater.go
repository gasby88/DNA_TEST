package testcase

//import (
//	"DNA_TEST/testframework"
//	"time"
//)
//
//func TestStateUpdaterTransction(ctx *testframework.TestFrameworkContext) bool {
//	k, _  := ctx.DnaClient.Account1.PublicKey.EncodePoint(true)
//	ctx.LogInfo("%x", k)
//	namespace := []byte("bing")
//	tx, err := ctx.Dna.NewStateUpdaterTransaction(ctx.DnaClient.Account1, false, namespace)
//	if err != nil {
//		ctx.LogError("NewStateUpdaterTransaction error:%s", err)
//		return false
//	}
//	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, tx)
//	if err != nil {
//		ctx.LogError("SendTransaction error:%s", err)
//		return false
//	}
//	_, err = ctx.Dna.WaitForGenerateBlock(30*time.Second, 1)
//	if err != nil {
//		ctx.LogError("WaitForGenerateBlock error:%s", err)
//		return false
//	}
//	ctx.LogInfo("Step 1")
//
//	key := []byte("Hello")
//	value := []byte("Word")
//	tx, err = ctx.Dna.NewStateUpdateTransction(ctx.DnaClient.Account1, namespace, key, value)
//	if err != nil {
//		ctx.LogError("NewStateUpdateTransction error:%s", err)
//		return false
//	}
//	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, tx)
//	if err == nil {
//		ctx.LogError("StateUpdateTransction should failed. Because of havn't namespace:%s", namespace)
//		return false
//	}
//
//	ctx.LogInfo("Step 2")
//	tx, err = ctx.Dna.NewStateUpdaterTransaction(ctx.DnaClient.Account1, true, namespace)
//	if err != nil {
//		ctx.LogError("NewStateUpdaterTransaction error:%s", err)
//		return false
//	}
//	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, tx)
//	if err != nil {
//		ctx.LogError("SendTransaction error:%s", err)
//		return false
//	}
//	_, err = ctx.Dna.WaitForGenerateBlock(30*time.Second, 1)
//	if err != nil {
//		ctx.LogError("WaitForGenerateBlock error:%s", err)
//		return false
//	}
//	ctx.LogInfo("Step 3")
//	tx, err = ctx.Dna.NewStateUpdateTransction(ctx.DnaClient.Account1, []byte(namespace), key, value)
//	if err != nil {
//		ctx.LogError("NewStateUpdateTransction error:%s", err)
//		return false
//	}
//	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, tx)
//	if err != nil {
//		ctx.LogError("StateUpdateTransction error:%s", err)
//		return false
//	}
//	ctx.LogInfo("Step 4")
//	return true
//}
