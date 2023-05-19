package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/itering/subscan/util"
	"github.com/stretchr/testify/assert"
)

func TestDao_SetHeartBeatNow(t *testing.T) {
	ctx := context.TODO()
	err := testDao.SetHeartBeatNow(ctx, "testAction")
	assert.NoError(t, err)
	assert.Equal(t, testDao.getCacheInt64(ctx, "testAction"), time.Now().Unix())
}

func TestDao_DaemonHealth(t *testing.T) {
	ctx := context.TODO()
	assert.Equal(t, testDao.DaemonHealth(ctx), map[string]bool{"substrate": false})
	err := testDao.SetHeartBeatNow(ctx, fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, "substrate"))
	assert.NoError(t, err)
	assert.Equal(t, testDao.DaemonHealth(ctx), map[string]bool{"substrate": true})
}
