package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_FillAlreadyBlockNum(t *testing.T) {
	ctx := context.TODO()
	err := testDao.SaveFillAlreadyBlockNum(ctx, 900000)
	assert.NoError(t, err)
	num, err := testDao.GetFillBestBlockNum(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 900000, num)
}

func TestDao_GetFillFinalizedBlockNum(t *testing.T) {
	ctx := context.TODO()
	err := testDao.SaveFillAlreadyFinalizedBlockNum(ctx, 899999)
	assert.NoError(t, err)
	num, err := testDao.GetFillFinalizedBlockNum(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 899999, num)
}

func TestDao_GetBlockByHash(t *testing.T) {
	ctx := context.TODO()
	block := testDao.GetBlockByHash(ctx, "0xd68b38c412404a4b5d4974e6dbb4a491ed7b6200d4edc24152693804441ce99d")
	assert.Equal(t, testBlock.BlockNum, block.BlockNum)
}

func TestDao_GetBlockByNum(t *testing.T) {
	block := testDao.GetBlockByNum(947687)
	assert.Equal(t, testBlock.ParentHash, block.ParentHash)
}

func TestDao_GetBlockList(t *testing.T) {
	blocks := testDao.GetBlockList(0, 100)
	assert.GreaterOrEqual(t, 1, len(blocks))
}

func TestDao_UpdateEventAndExtrinsic(t *testing.T) {
	txn := testDao.DbBegin()
	err := testDao.UpdateEventAndExtrinsic(txn, &testBlock, 1, 2, 1594791900, "60e2feb892e672d5579ed10ecae0d162031fe5adc3692498ad262fb126a65732", false, true)
	assert.NoError(t, err)
	txn.Commit()
	block := testDao.GetBlockByNum(947687)
	assert.Equal(t, 1594791900, block.BlockTimestamp)
}

func TestDao_GetNearBlock(t *testing.T) {
	txn := testDao.DbBegin()
	_ = testDao.CreateBlock(txn, &model.ChainBlock{BlockNum: 947688, Hash: "0xd68b38c412404a4b5d4974e6dbb4a491ed7b6200d4edc24152693804441ce999"})
	txn.Commit()
	block := testDao.GetNearBlock(947687)
	assert.Equal(t, 947688, block.BlockNum)
}

func TestDaoSetBlockFinalized(t *testing.T) {
	block := testDao.GetBlockByNum(947687)
	testDao.SetBlockFinalized(block)
	block = testDao.GetBlockByNum(947687)
	assert.Equal(t, true, block.Finalized)
}

func TestDao_BlocksReverseByNum(t *testing.T) {
	blockMap := testDao.BlocksReverseByNum([]int{947687})
	assert.Equal(t, map[int]model.ChainBlock{947687: *testDao.GetBlockByNum(947687)}, blockMap)
}
