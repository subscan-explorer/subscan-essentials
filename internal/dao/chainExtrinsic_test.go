package dao

import (
	"context"
	"encoding/json"
	"github.com/itering/subscan/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_CreateExtrinsic(t *testing.T) {
	ctx := context.TODO()
	txn := testDao.DbBegin()
	_ = testDao.CreateExtrinsic(ctx, txn, &testExtrinsic)
	txn.Commit()
}

func TestDao_DropExtrinsicNotFinalizedData(t *testing.T) {
	ctx := context.TODO()
	txn := testDao.DbBegin()

	tempExtrinsic := testExtrinsic
	tempExtrinsic.BlockNum = 947688
	err := testDao.CreateExtrinsic(ctx, txn, &tempExtrinsic)
	txn.Commit()
	assert.NoError(t, err)

	testDao.DropExtrinsicNotFinalizedData(ctx, tempExtrinsic.BlockNum, true)
	assert.Equal(t, []model.ChainExtrinsicJson(nil), testDao.GetExtrinsicsByBlockNum(947688))
}

func TestDao_GetExtrinsicsByBlockNum(t *testing.T) {
	extrinsics := testDao.GetExtrinsicsByBlockNum(947687)
	assert.Equal(t, []model.ChainExtrinsicJson{{BlockTimestamp: 1594791900, BlockNum: 947687, ExtrinsicIndex: "947687-0", CallModuleFunction: "set", CallModule: "timestamp", Params: "null", AccountId: "", AccountIndex: "", Signature: "", Nonce: 0, ExtrinsicHash: "", Success: true, Fee: decimal.New(0, 0)}}, extrinsics)
}

func TestDao_GetExtrinsicsByHash(t *testing.T) {
	ctx := context.TODO()
	extrinsics := testDao.GetExtrinsicsByHash(ctx, "0x368f61800f8645f67d59baf0602b236ff47952097dcaef3aa026b50ddc8dcea0")
	expect := testSignedExtrinsic
	params, _ := json.Marshal(testSignedExtrinsic.Params)
	expect.Params = []byte(string(params))
	expect.Fee = decimal.Zero
	extrinsics.Fee = decimal.Zero
	assert.EqualValues(t, &expect, extrinsics)
}

func TestDao_ExtrinsicList(t *testing.T) {
	ctx := context.TODO()
	extrinsic, _ := testDao.GetExtrinsicList(ctx, 0, 100, "desc")
	assert.GreaterOrEqual(t, 2, len(extrinsic))
}

func TestDao_GetExtrinsicsDetailByIndex(t *testing.T) {
	ctx := context.TODO()
	extrinsic := testDao.GetExtrinsicsDetailByIndex(ctx, "947689-1")
	assert.Equal(t, "5Ct9YwxrnJmMw1q45iD73u55XJcM6cyfwnVyNrHJQCC52zia", extrinsic.AccountId)
	assert.Equal(t, testSignedExtrinsic.Params, extrinsic.Params)
}

func TestDao_ExtrinsicsAsJson(t *testing.T) {
	ctx := context.TODO()
	extrinsics := testDao.GetExtrinsicsByHash(ctx, "0x368f61800f8645f67d59baf0602b236ff47952097dcaef3aa026b50ddc8dcea0")
	assert.Equal(t, `[{"name":"dest","type":"Address","value":"563d11af91b3a166d07110bb49e84094f38364ef39c43a26066ca123a8b9532b","valueRaw":""},{"name":"value","type":"Compact\u003cBalance\u003e","value":"1000000000000000000","valueRaw":""}]`, testDao.ExtrinsicsAsJson(extrinsics).Params)
}
