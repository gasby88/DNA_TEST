package main

import (
	"DNA/common/log"
	."DNA_TEST/dna"
	."DNA_TEST/testframework"
	_"DNA_TEST/testcase"
	"flag"
	log4 "github.com/alecthomas/log4go"
	"time"
)

var DNAJsonRpcAddress string

func init() {
	flag.StringVar(&DNAJsonRpcAddress, "rpc", "http://localhost:20336", "The address of dna jsonrpc")
	flag.Parse()
}

func main() {
	log4.LoadConfiguration("./etc/log4go.xml")
	log.CreatePrintLog("./log")

	dna := NewDna([]string{DNAJsonRpcAddress})
	dnaClient := NewDnaClient()
	dnaClient.Init()

	TFramework.SetDna(dna)
	TFramework.SetDnaClient(dnaClient)
	TFramework.Start()

	time.Sleep(time.Second)
}
