package service

import (
	"github.com/itering/subscan/model"
	"testing"
)

func Test_emitEvent(t *testing.T) {
	testSrv.emitEvent(&testBlock, &testEvent)
}

func Test_emitExtrinsic(t *testing.T) {
	testSrv.emitExtrinsic(&testBlock, &testSignedExtrinsic, []model.ChainEvent{testEvent})
}
