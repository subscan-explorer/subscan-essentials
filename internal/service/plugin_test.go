package service

import (
	"context"
	"github.com/itering/subscan/model"
	"testing"
)

func Test_emitEvent(t *testing.T) {
	testSrv.emitEvent(&testBlock, &testEvent)
}

func Test_emitExtrinsic(t *testing.T) {
	testSrv.emitExtrinsic(context.TODO(), &testBlock, &testSignedExtrinsic, []model.ChainEvent{testEvent})
}
