package dao

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDao_SetHeartBeatNow(t *testing.T) {
	ctx := context.TODO()
	err := testDao.SetHeartBeatNow(ctx, "testAction")
	assert.NoError(t, err)
	assert.Equal(t, testDao.getCacheInt64(ctx, "testAction"), time.Now().Unix())

}
