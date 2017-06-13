package dna

import (
	"DNA/client"
	"fmt"
)

type DnaClient struct {
	Client   client.Client
	Admin    *client.Account
	Account1 *client.Account
	Account2 *client.Account
	Account3 *client.Account
	Account4 *client.Account
	Account5 *client.Account
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

func (this *DnaClient) CreateClient(name string) *client.ClientImpl {
	path := fmt.Sprintf("./wallet_%s.txt", name)
	if FileExisted(path) {
		return client.OpenClient(path, []byte("dna"))
	} else {
		return client.CreateClient(path, []byte("dna"))
	}
}
