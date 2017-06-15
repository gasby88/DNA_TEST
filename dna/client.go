package dna

import (
	"DNA/account"
	"DNA/common"
	"fmt"
)

type DnaClient struct {
	Client   account.Client
	Admin    *account.Account
	Account1 *account.Account
	Account2 *account.Account
	Account3 *account.Account
	Account4 *account.Account
	Account5 *account.Account
}

func NewDnaClient() *DnaClient {
	return &DnaClient{}
}

func (this *DnaClient)Init(){
	testClient := this.CreateClient("test")
	this.Client = testClient

	var err error
	this.Admin, err = testClient.CreateAccount()
	if err != nil {
		panic(fmt.Errorf("DnaClient CreateAccount:Admin error:%s", err))
	}
	this.Account1, err = testClient.CreateAccount()
	if err != nil {
		panic(fmt.Errorf("DnaClient CreateAccount:Account1 error:%s", err))
	}
	this.Account2, err = testClient.CreateAccount()
	if err != nil {
		panic(fmt.Errorf("DnaClient CreateAccount:Account2 error:%s", err))
	}
	this.Account3, err = testClient.CreateAccount()
	if err != nil {
		panic(fmt.Errorf("DnaClient CreateAccount:Account3 error:%s", err))
	}
	this.Account4, err = testClient.CreateAccount()
	if err != nil {
		panic(fmt.Errorf("DnaClient CreateAccount:Account4 error:%s", err))
	}
	this.Account5, err = testClient.CreateAccount()
	if err != nil {
		panic(fmt.Errorf("DnaClient CreateAccount:Account5 error:%s", err))
	}
}

func (this *DnaClient) CreateClient(name string) *account.ClientImpl {
	path := fmt.Sprintf("./wallet_%s.txt", name)
	if common.FileExisted(path) {
		return account.Open(path, []byte("dna"))
	} else {
		return account.Create(path, []byte("dna"))
	}
}
