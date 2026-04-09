package dao

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
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

func TestDao_GetExtrinsicCountUsesHistoricalCacheButLatestDB(t *testing.T) {
	ctx := context.TODO()
	seedLatestExtrinsicFixtures(t)

	err := testDao.SaveFillAlreadyBlockNum(ctx, int(model.SplitTableBlockNum+1))
	assert.NoError(t, err)

	conn, _ := testDao.redis.Redis().GetContext(ctx)
	defer conn.Close()

	historyAllKey := testDao.extrinsicCountCacheKey(0)
	historySignedKey := testDao.extrinsicCountCacheKey(0, model.Where("is_signed = ?", true))
	latestAllKey := testDao.extrinsicCountCacheKey(1)
	latestSignedKey := testDao.extrinsicCountCacheKey(1, model.Where("is_signed = ?", true))

	_, err = conn.Do("SET", historyAllKey, 100)
	assert.NoError(t, err)
	_, err = conn.Do("SET", historySignedKey, 7)
	assert.NoError(t, err)
	_, err = conn.Do("SET", latestAllKey, 200)
	assert.NoError(t, err)
	_, err = conn.Do("SET", latestSignedKey, 300)
	assert.NoError(t, err)

	assert.EqualValues(t, 101, testDao.GetExtrinsicCount(ctx))
	assert.EqualValues(t, 8, testDao.GetExtrinsicCount(ctx, model.Where("is_signed = ?", true)))
}

func TestDao_GetExtrinsicCountBackfillsHistoricalCacheWithJitteredTTL(t *testing.T) {
	ctx := context.TODO()
	seedLatestExtrinsicFixtures(t)

	err := testDao.SaveFillAlreadyBlockNum(ctx, int(model.SplitTableBlockNum+1))
	assert.NoError(t, err)

	conn, _ := testDao.redis.Redis().GetContext(ctx)
	defer conn.Close()

	historyAllKey := testDao.extrinsicCountCacheKey(0)
	_, err = conn.Do("DEL", historyAllKey)
	assert.NoError(t, err)

	assert.EqualValues(t, 3, testDao.GetExtrinsicCount(ctx))

	cachedCount, err := redis.Int64(conn.Do("GET", historyAllKey))
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cachedCount)

	ttl, err := redis.Int(conn.Do("TTL", historyAllKey))
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, ttl, extrinsicCountCacheTTL-1)
	assert.LessOrEqual(t, ttl, extrinsicCountCacheTTL+extrinsicCountCacheTTLJitter)
}

func seedLatestExtrinsicFixtures(t *testing.T) {
	t.Helper()

	testDao.AddIndex(model.SplitTableBlockNum)

	ctx := context.TODO()
	txn := testDao.DbBegin()
	err := testDao.CreateExtrinsic(ctx, txn, []model.ChainExtrinsic{{
		ID:                 model.SplitTableBlockNum*model.IdGenerateCoefficient + 1,
		ExtrinsicIndex:     fmt.Sprintf("%d-1", model.SplitTableBlockNum+1),
		BlockNum:           model.SplitTableBlockNum + 1,
		BlockTimestamp:     1594791901,
		CallModuleFunction: "remark",
		CallModule:         "system",
		Success:            true,
		ExtrinsicHash:      "0xsplit-table-extrinsic",
		IsSigned:           true,
	}}, 1)
	assert.NoError(t, err)
	assert.NoError(t, txn.Commit().Error)
}
