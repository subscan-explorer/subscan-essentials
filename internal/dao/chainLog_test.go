package dao

import (
	"testing"

	"github.com/itering/subscan/model"
	"github.com/stretchr/testify/assert"
)

func TestDao_CreateLog(t *testing.T) {
	txn := testDao.DbBegin()
	err := testDao.CreateLog(txn, &testLog)
	txn.Commit()
	assert.NoError(t, err)
}

func TestDao_DropLogsNotFinalizedData(t *testing.T) {
	txn := testDao.DbBegin()
	testLog.BlockNum = 947688
	txn.Commit()
	testDao.DropLogsNotFinalizedData(947688, true)
	assert.Equal(t, []model.ChainLogJson(nil), testDao.GetLogByBlockNum(947688))
}

func TestDao_GetLogsByIndex(t *testing.T) {
	log := testDao.GetLogsByIndex("947687-0")
	assert.Equal(t, 947687, log.BlockNum)
}

func TestDao_GetLogByBlockNum(t *testing.T) {
	logs := testDao.GetLogByBlockNum(947687)
	for _, log := range logs {
		assert.Equal(t, 947687, log.BlockNum)
	}
}
