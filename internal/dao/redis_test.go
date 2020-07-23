package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/util"
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

func TestDao_DaemonHeath(t *testing.T) {
	ctx := context.TODO()
	assert.Equal(t, testDao.DaemonHeath(ctx), map[string]bool{"substrate": false})
	err := testDao.SetHeartBeatNow(ctx, fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, "substrate"))
	assert.NoError(t, err)
	assert.Equal(t, testDao.DaemonHeath(ctx), map[string]bool{"substrate": true})
}
