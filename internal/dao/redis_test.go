package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDao_SetHeartBeatNow(t *testing.T) {
	ctx := context.TODO()
	err := TestDao.SetHeartBeatNow(ctx, "testAction")
	assert.NoError(t, err)
	assert.Equal(t, TestDao.GetCacheInt64(ctx, "testAction"), time.Now().Unix())

}

func TestDao_DaemonHeath(t *testing.T) {
	ctx := context.TODO()
	assert.Equal(t, TestDao.DaemonHeath(ctx), map[string]bool{"substrate": false})
	err := TestDao.SetHeartBeatNow(ctx, fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, "substrate"))
	assert.NoError(t, err)
	assert.Equal(t, TestDao.DaemonHeath(ctx), map[string]bool{"substrate": true})
}
