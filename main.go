package main

import (
	"DNA/common/log"
	. "DNA_TEST/dna"
	_ "DNA_TEST/testcase"
	. "DNA_TEST/testframework"
	"flag"
	log4 "github.com/alecthomas/log4go"
	"strings"
	"time"
)

var (
	DNAJsonRpcAddress string
	DNAWSAddress      string
	CycleTestMode     bool
	CycleTestInterval int
	BenchTestMode     bool
	BenchThreadNum    int
	BenchLastTime     int
)

func init() {
	flag.StringVar(&DNAJsonRpcAddress, "rpc", "http://localhost:20336", "The address of dna jsonrpc")
	flag.StringVar(&DNAWSAddress, "ws", "ws://localhost:20335", "The address of dna web socket")
	flag.BoolVar(&CycleTestMode, "c", false, "Is cycle test mode")
	flag.IntVar(&CycleTestInterval, "ci", 10, "Interval between test in cycle mode")
	flag.BoolVar(&BenchTestMode, "b", false, "Is benchmark test mode")
	flag.IntVar(&BenchThreadNum, "bn", 1, "Thread num in benchmark mode")
	flag.IntVar(&BenchLastTime, "bt", 30, "Bench time in benchmark mode")
	flag.Parse()
}

func parseAddress(addresses string) []string {
	return strings.Split(strings.Trim(addresses, ";"), ";")
}

func main() {
	log4.LoadConfiguration("./etc/log4go.xml")
	log.Init("./log", log.Stdout)

	dna := NewDna(parseAddress(DNAJsonRpcAddress), parseAddress(DNAWSAddress))
	dnaClient := NewDnaClient()
	dnaClient.Init()

	TFramework.SetDna(dna)
	TFramework.SetDnaClient(dnaClient)
	TFramework.Start(&TestFrameworkOptions{
		CycleTestMode:     CycleTestMode,
		CycleTestInterval: CycleTestInterval,
		BenchTestMode:     BenchTestMode,
		BenchThreadNum:    BenchThreadNum,
		BenchLastTime:     BenchLastTime,
	})

	time.Sleep(time.Second)
}
