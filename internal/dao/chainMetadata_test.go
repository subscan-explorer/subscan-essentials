package dao

import (
	"context"
	"github.com/garyburd/redigo/redis"
	"github.com/itering/subscan/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_SetMetadata(t *testing.T) {
	ctx := context.TODO()
	testCase := map[string]interface{}{"key1": "value1", "key2": "value2"}

	err := testDao.SetMetadata(ctx, testCase)
	assert.NoError(t, err)

	conn := testDao.redis.Get(ctx)
	defer conn.Close()
	for key, expect := range testCase {
		value, _ := redis.String(conn.Do("HGET", RedisMetadataKey, key))
		assert.Equal(t, expect, value)
	}
}

func TestDao_IncrMetadata(t *testing.T) {
	ctx := context.TODO()
	testCase := map[string]interface{}{"key1": "1", "key2": "2"}
	_ = testDao.SetMetadata(ctx, testCase)

	conn := testDao.redis.Get(ctx)
	defer conn.Close()

	for key, expect := range testCase {
		_ = testDao.IncrMetadata(ctx, key, 1)
		value, _ := redis.Int(conn.Do("HGET", RedisMetadataKey, key))
		assert.Equal(t, util.StringToInt(expect.(string))+1, value)
	}
}

func TestDao_GetMetadata(t *testing.T) {
	ctx := context.TODO()
	testCase := map[string]interface{}{"key1": "1", "key2": "2"}
	_ = testDao.SetMetadata(ctx, testCase)
	metadata, err := testDao.GetMetadata(ctx)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"key1": "1", "key2": "2"}, metadata)
}

func TestDao_GetCurrentBlockNum(t *testing.T) {
	ctx := context.TODO()
	testCase := map[string]interface{}{"blockNum": 999999}
	_ = testDao.SetMetadata(ctx, testCase)
	blockNum, err := testDao.GetBestBlockNum(ctx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(999999), blockNum)
}

func TestDao_GetFinalizedBlockNum(t *testing.T) {
	ctx := context.TODO()
	testCase := map[string]interface{}{"finalized_blockNum": 999999}
	_ = testDao.SetMetadata(ctx, testCase)
	blockNum, err := testDao.GetFinalizedBlockNum(ctx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(999999), blockNum)
}
