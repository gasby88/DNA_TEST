package testframework

import (
	"DNA_TEST/dna"
	log4 "github.com/alecthomas/log4go"
)

type TestFrameworkContext struct {
	Dna       *dna.Dna
	DnaClient *dna.DnaClient
	DnaAsset  *dna.DnaAsset
	failNowCh chan interface{}
}

func NewTestFrameworkContext(dna *dna.Dna,
	dnaClient *dna.DnaClient,
	dnaAsset *dna.DnaAsset,
	failNowCh chan interface{}) *TestFrameworkContext {
	return &TestFrameworkContext{
		Dna:       dna,
		DnaClient: dnaClient,
		DnaAsset:  dnaAsset,
		failNowCh: failNowCh,
	}
}

func (this *TestFrameworkContext) LogInfo(arg0 interface{}, args ...interface{}) {
	log4.Info(arg0, args...)
}

func (this *TestFrameworkContext) LogError(arg0 interface{}, args ...interface{}) {
	log4.Error(arg0, args...)
}

func (this *TestFrameworkContext) FailNow() {
	select {
	case <-this.failNowCh:
	default:
		close(this.failNowCh)
	}
}
