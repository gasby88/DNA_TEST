package testcase
//
//import (
//	"DNA_TEST/testframework"
//	"encoding/json"
//	"fmt"
//	"math/rand"
//	"ontsdk/key"
//	"ontsdk/ont"
//	"ontsdk/ont/claimcnt"
//	"ontsdk/ont/id/did"
//	"time"
//)
//
//var (
//	OntId1    string
//	OntDDO1   *ont.DDO
//	OntClaim1 *ont.VerifiableClaim
//	OntOwner1 *ont.Owner
//	OntId2    string
//	OntDDO2   *ont.DDO
//	OntClaim2 *ont.VerifiableClaim
//	OntOwner2 *ont.Owner
//)
//
//func init() {
//	rand.Seed(time.Now().UnixNano())
//
//	OntId1 = did.NewDID(fmt.Sprintf("%x", rand.Int())).String()
//	ecdsaKey1, err := key.NewECDSACryptKey()
//	if err != nil {
//		fmt.Printf("NewECDSACryptKey error:%s\n", err)
//		return
//	}
//	OntDDO1 = ont.NewDDO(OntId1)
//	OntOwner1, err = OntDDO1.AddOwner(ecdsaKey1, time.Now().AddDate(0, 1, 0))
//	if err != nil {
//		fmt.Printf("OntDDO1 AddOwner error:%s\n", err)
//		return
//	}
//	OntDDO1.SetUpdateTime()
//	err = OntDDO1.SignatureWithDefaultOwner()
//	if err != nil {
//		fmt.Printf("OntDDO1.SignatureWithDefaultOwner error:%s\n", err)
//		return
//	}
//	OntClaim1 = ont.NewVerifiableClaim("OntClaim1", claimcnt.ONT_IDENTITY_CLAIM, OntOwner1)
//	//OntClaim1.SetExpires(time.Now().AddDate(1, 0, 0))
//
//	OntId2 = did.NewDID(fmt.Sprintf("%x", rand.Int())).String()
//	ecdsaKey2, err := key.NewECDSACryptKey()
//	if err != nil {
//		fmt.Printf("NewECDSACryptKey error:%s\n", err)
//		return
//	}
//	OntDDO2 = ont.NewDDO(OntId2)
//	OntOwner2, err = OntDDO2.AddOwner(ecdsaKey2, time.Now().AddDate(0, 1, 0))
//	if err != nil {
//		fmt.Printf("OntDDO1 AddOwner error:%s\n", err)
//		return
//	}
//	OntDDO2.SetUpdateTime()
//	err = OntDDO2.SignatureWithDefaultOwner()
//	if err != nil {
//		fmt.Printf("OntDDO2.SignatureWithDefaultOwner error:%s\n", err)
//		return
//	}
//	OntClaim2 = ont.NewVerifiableClaim("OntClaim2", claimcnt.ONT_IDENTITY_CLAIM, OntOwner2)
//	OntClaim2.SetExpires(time.Now().AddDate(1, 0, 0))
//	//OntClaim2.SetContent(&claimcnt.IdentityClaimContent{
//	//	OntId:OntId2,
//	//	IdentityType:claimcnt.ONT_IDENTITY_TYPE_CA,
//	//	IdentityData:[]byte("Hello world"),
//	//})
//	//err = OntClaim2.Signature()
//	//if err != nil {
//	//	fmt.Printf("OntClaim1 Signature error:%s\n", err)
//	//	return
//	//}
//}
//
//func TestOntIdentiyUpdate(ctx *testframework.TestFrameworkContext) bool {
//	data, err := json.Marshal(OntDDO1)
//	if err != nil {
//		ctx.LogError("TestOntIdentiyUpdate json.Marshal error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//	tx, err := ctx.Dna.NewIdentityUpdateTransaction([]byte(OntId1), data, ctx.DnaClient.Account1.PublicKey)
//	if err != nil {
//		ctx.LogError("TestOntIdentiyUpdate NewIdentityUpdateTransaction error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, tx)
//	if err != nil {
//		ctx.LogError("TestOntIdentiyUpdate SendTransaction IdentityUpdateTransaction error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	_, err = ctx.Dna.WaitForGenerateBlock(30 * time.Second)
//	if err != nil {
//		ctx.LogError("TestOntIdentiyUpdate WaitForGenerateBlock error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	ddo, err := ctx.Dna.GetIdentity(OntId1)
//	if err != nil {
//		ctx.LogError("GetIdentity error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	if ddo == nil {
//		ctx.LogError("GetIdentity by ontId:%s nil", err, OntId1)
//		ctx.FailNow()
//		return false
//	}
//
//	if ddo.OntId != OntId1 {
//		ctx.LogError("GetIdentity error:%s OntId:%s != %s", ddo.OntId, OntId1)
//		ctx.FailNow()
//		return false
//	}
//
//	tx2, err := ctx.Dna.NewIdentityUpdateTransaction([]byte(OntId1), data, ctx.DnaClient.Account1.PublicKey)
//	if err != nil {
//		ctx.LogError("TestOntIdentiyUpdate NewIdentityUpdateTransaction2 error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, tx2)
//	if err != nil {
//		ctx.LogError("TestOntIdentiyUpdate SendTransaction IdentityUpdateTransaction error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	return true
//}
//
//func TestOntIdentiyClaimUpdate(ctx *testframework.TestFrameworkContext) bool {
//	nonce := fmt.Sprintf("%d", rand.Int())
//
//	idClaimCnt := claimcnt.NewIdentityClaimContent(OntId1).AddIdentityItem(claimcnt.NewIdentityClaimContenttItem(
//		"shca",
//		claimcnt.ONT_IDENTITY_TYPE_CA,
//		[]byte("Hello world"),
//		nonce,
//	))
//	OntClaim1.SetContent(idClaimCnt)
//	err := OntClaim1.Signature()
//	if err != nil {
//		ctx.LogError("OntClaim1 Signature error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	data, err := json.Marshal(OntClaim1)
//	if err != nil {
//		ctx.LogError("json.Marshal OntClaim1 error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	ctx.LogInfo("OntClaim1:%s", data)
//	tx, err := ctx.Dna.NewIdentityClaimUpdateTransaction([]byte(OntId1), data, ctx.DnaClient.Account1.PublicKey)
//	if err != nil {
//		ctx.LogError("NewIdentityClaimUpdateTransaction error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	_, err = ctx.Dna.SendTransaction(ctx.DnaClient.Account1, tx)
//	if err != nil {
//		ctx.LogError("SendTransaction error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	_, err = ctx.Dna.WaitForGenerateBlock(30 * time.Second)
//	if err != nil {
//		ctx.LogError("TestOntIdentiyUpdate WaitForGenerateBlock error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	idClaim, err := ctx.Dna.GetIdentityClaim(OntId1)
//	if err != nil {
//		ctx.LogError("GetIdentityClaim error:%s", err)
//		ctx.FailNow()
//		return false
//	}
//
//	if idClaim == nil {
//		ctx.LogError("GetIdentityClaim by OntId:%s nil", OntId1)
//		ctx.FailNow()
//		return false
//	}
//
//	if idClaim.Type[1] != claimcnt.ONT_IDENTITY_CLAIM{
//		ctx.LogError("IdentityClaimContent claim type:%s not:%s", idClaim.Type[1], claimcnt.ONT_IDENTITY_CLAIM)
//		ctx.FailNow()
//		return false
//	}
//	identityClaim := &claimcnt.IdentityClaimContent{}
//	err = json.Unmarshal(idClaim.Content, identityClaim)
//	if err != nil {
//		ctx.LogError("IdentityClaimContent json.Unmarshal IdentityClaimContent error:%s", idClaim.Content)
//		ctx.FailNow()
//		return false
//	}
//
//	if identityClaim.Identities[0].Desc != nonce {
//		ctx.LogError("IdentityClaimContent Data error. Nonce:%s != %s", identityClaim.Identities[0].Desc, nonce)
//		ctx.FailNow()
//		return false
//	}
//
//	return true
//}
