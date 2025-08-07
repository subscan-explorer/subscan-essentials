package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_emitEvent(t *testing.T) {
	assert.NoError(t, testSrv.emitEvent(&testEvent))
}

func Test_emitExtrinsic(t *testing.T) {
	assert.NoError(t, testSrv.emitExtrinsic(context.TODO(), &testSignedExtrinsic))
}
