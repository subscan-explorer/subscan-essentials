package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_CreateExtrinsic(t *testing.T) {
	ctx := context.TODO()
	txn := testDao.DbBegin()
	_ = testDao.CreateExtrinsic(ctx, txn, []model.ChainExtrinsic{testExtrinsic}, 1)
	txn.Commit()
}

func TestDao_GetExtrinsicsByHash(t *testing.T) {
	ctx := context.TODO()
	extrinsics := testDao.GetExtrinsicsByHash(ctx, "0x368f61800f8645f67d59baf0602b236ff47952097dcaef3aa026b50ddc8dcea0")
	expect := testSignedExtrinsic
	expect.Params = testSignedExtrinsic.Params
	expect.Fee = decimal.Zero
	expect.UsedFee = decimal.Zero
	extrinsics.Fee = decimal.Zero
	extrinsics.UsedFee = decimal.Zero
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
	assert.Equal(t, "1pShHDve62qNYqa3MG7C3uENvbznvXp2HETY9GexHDbDUyC", extrinsic.AccountId)
	assert.Equal(t, testSignedExtrinsic.Params, extrinsic.Params)
}

func TestDao_ExtrinsicsAsJson(t *testing.T) {
	ctx := context.TODO()
	extrinsics := testDao.GetExtrinsicsByHash(ctx, "0x368f61800f8645f67d59baf0602b236ff47952097dcaef3aa026b50ddc8dcea0")
	assert.Equal(t, `[{"name":"dest","type":"Address","value":"563d11af91b3a166d07110bb49e84094f38364ef39c43a26066ca123a8b9532b"},{"name":"value","type":"Compact\u003cBalance\u003e","value":"1000000000000000000"}]`, util.ToString(testDao.ExtrinsicsAsJson(extrinsics).Params))
}
