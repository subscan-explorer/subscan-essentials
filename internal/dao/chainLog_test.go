package dao

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_CreateLog(t *testing.T) {
	txn := testDao.DbBegin()
	err := testDao.CreateLog(txn, &testLog)
	txn.Commit()
	assert.NoError(t, err)

}

func TestDao_GetLogByBlockNum(t *testing.T) {
	logs := testDao.GetLogByBlockNum(context.TODO(), 947687)
	for _, log := range logs {
		assert.Equal(t, 947687, log.BlockNum)
	}
}
