package testframework

import (
	. "DNA_TEST/dna"
	"fmt"
	log4 "github.com/alecthomas/log4go"
	"reflect"
	"time"
)

var TFramework = NewTestFramework()

type TestCase func(ctx *TestFrameworkContext) bool

type TestFrameworkOptions struct {
	CycleTestMode     bool
	CycleTestInterval int
	BenchTestMode     bool
	BenchThreadNum    int
	BenchLastTime     int
}

type TestFramework struct {
	options        *TestFrameworkOptions
	startTime      time.Time
	testCases      []TestCase
	benchTestCases []TestCase
	testCasesMap   map[string]string
	testCaseRes    map[string]bool
	dna            *Dna
	dnaClient      *DnaClient
	dnaAsset       *DnaAsset
}

func NewTestFramework() *TestFramework {
	return &TestFramework{
		testCases:    make([]TestCase, 0),
		testCasesMap: make(map[string]string, 0),
		testCaseRes:  make(map[string]bool, 0),
		dnaAsset:     NewDnaAsset(),
	}
}

func (this *TestFramework) RegTestCase(name string, testCase TestCase) {
	this.testCases = append(this.testCases, testCase)
	this.testCasesMap[this.getTestCaseId(testCase)] = name
}

func (this *TestFramework) RegBenchTestCase(name string, testCase TestCase) {
	this.benchTestCases = append(this.testCases, testCase)
	this.testCasesMap[this.getTestCaseId(testCase)] = name
}

func (this *TestFramework) Start(ops *TestFrameworkOptions) {
	this.options = ops
	if this.options.BenchTestMode {
		go this.runTestList(this.benchTestCases)
	}

	if !this.options.CycleTestMode {
		this.runTestList(this.testCases)
		return
	}

	timer := time.NewTimer(time.Duration(this.options.CycleTestInterval) * time.Second)
	for {
		select {
		case <-timer.C:
			this.runTestList(this.testCases)
			timer.Reset(time.Duration(this.options.CycleTestInterval) * time.Second)
		}
	}
}

func (this *TestFramework) runTestList(testCaseList []TestCase) {
	this.onTestStart()
	defer this.onTestFinish()
	failNowCh := make(chan interface{}, 0)
	for i, testCase := range this.testCases {
		select {
		case <-failNowCh:
			this.onTestFailNow()
			return
		default:
			this.runTest(i+1, failNowCh, testCase)
		}
	}
}

func (this *TestFramework) runTest(index int, failNowCh chan interface{}, testCase TestCase) {
	ctx := NewTestFrameworkContext(this.dna,
		this.dnaClient,
		this.dnaAsset,
		failNowCh,
		this.options.BenchThreadNum,
		time.Duration(this.options.BenchLastTime)*time.Second)
	this.onBeforeTestCaseStart(index, testCase)
	ok := testCase(ctx)
	this.onAfterTestCaseFinish(index, testCase, ok)
	this.testCaseRes[this.getTestCaseId(testCase)] = ok
}

func (this *TestFramework) SetDna(dna *Dna) {
	this.dna = dna
}

func (this *TestFramework) SetDnaClient(dnaClient *DnaClient) {
	this.dnaClient = dnaClient
}

func (this *TestFramework) onTestStart() {
	version, _ := this.dna.GetVersion()

	log4.Info("\t\t\t===============================================================")
	log4.Info("\t\t\t-------DNA Test Start Version:%s", version)
	log4.Info("\t\t\t===============================================================")
	log4.Info("")
	this.startTime = time.Now()
}

func (this *TestFramework) onTestFinish() {

	failedList := make([]string, 0)
	successList := make([]string, 0)
	for testCase, ok := range this.testCaseRes {
		if ok {
			successList = append(successList, this.getTestCaseName(testCase))
		} else {
			failedList = append(failedList, this.getTestCaseName(testCase))
		}
	}

	skipList := make([]string, 0)
	for _, testCase := range this.testCases {
		_, ok := this.testCaseRes[this.getTestCaseId(testCase)]
		if !ok {
			skipList = append(skipList, this.getTestCaseName(testCase))
		}
	}

	succCount := len(successList)
	failedCount := len(failedList)

	log4.Info("\t\t===============================================================")
	log4.Info("\t\tDNA Test Finish Total:%v Success:%v Failed:%v Skip:%v TimeCost:%.2f s.",
		len(this.testCases),
		succCount,
		failedCount,
		len(this.testCases)-succCount-failedCount,
		time.Now().Sub(this.startTime).Seconds())
	if succCount > 0 {
		log4.Info("\t\t---------------------------------------------------------------")
		log4.Info("\t\t\tSuccess list:")
		for i, succCase := range successList {
			log4.Info("\t\t\t%d.\t%s", i+1, succCase)
		}
	}
	if failedCount > 0 {
		log4.Info("\t\t---------------------------------------------------------------")
		log4.Info("\t\t\tFail list:")
		for i, failCase := range failedList {
			log4.Info("\t\t\t%d.\t%s", i+1, failCase)
		}
	}
	if len(skipList) > 0 {
		log4.Info("\t\t---------------------------------------------------------------")
		log4.Info("\t\t\tSkip list:")
		for i, failCase := range skipList {
			log4.Info("\t\t\t%d.\t%s", i+1, failCase)
		}
	}
	log4.Info("\t\t===============================================================")
}

func (this *TestFramework) onTestFailNow() {
	log4.Info("Test Stop.")
}

func (this *TestFramework) onBeforeTestCaseStart(index int, testCase TestCase) {
	log4.Info("===============================================================")
	log4.Info("%d. Start TestCase:%s", index, this.getTestCaseName(testCase))
	log4.Info("---------------------------------------------------------------")
}

func (this *TestFramework) onAfterTestCaseFinish(index int, testCase TestCase, res bool) {
	if res {
		log4.Info("TestCase:%s success.", this.getTestCaseName(testCase))
	} else {
		log4.Info("TestCase:%s failed.", this.getTestCaseName(testCase))
	}
	log4.Info("---------------------------------------------------------------")
	log4.Info("")
}

func (this *TestFramework) getTestCaseName(testCase interface{}) string {
	testCaseStr, ok := testCase.(string)
	if !ok {
		testCaseStr = this.getTestCaseId(testCase)
	}
	name, ok := this.testCasesMap[testCaseStr]
	if ok {
		return name
	}
	return ""
}

func (this *TestFramework) getTestCaseId(testCase interface{}) string {
	return fmt.Sprintf("%v", reflect.ValueOf(testCase).Pointer())
}
