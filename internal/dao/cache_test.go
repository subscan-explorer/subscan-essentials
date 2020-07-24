package dao

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDao_SetCache(t *testing.T) {
	ctx := context.TODO()
	testCase := []struct {
		key    string
		value  interface{}
		expect string
	}{
		{"test1", "testvalue1", "testvalue1"},
		{"test2", int64(10), "10"},
		{"test3", 10, "10"},
		{"test4", map[string]string{"name": "amy"}, "{\"name\":\"amy\"}"},
		{"test5", "", ""},
	}
	for _, test := range testCase {
		err := testDao.setCache(ctx, test.key, test.value, 10)
		assert.NoError(t, err)
		assert.Equal(t, testDao.getCacheString(ctx, test.key), test.expect)
	}
}

func TestDao_GetCacheBytes(t *testing.T) {
	ctx := context.TODO()
	err := testDao.setCache(ctx, "test1", "value", 10)
	assert.NoError(t, err)
	assert.Equal(t, testDao.getCacheBytes(ctx, "test1"), []byte{0x76, 0x61, 0x6c, 0x75, 0x65})
}

func TestDao_GetCacheInt64(t *testing.T) {
	ctx := context.TODO()
	err := testDao.setCache(ctx, "test1", 2000, 10)
	assert.NoError(t, err)
	assert.Equal(t, testDao.getCacheInt64(ctx, "test1"), int64(2000))
}

func TestDao_GetCacheString(t *testing.T) {
	ctx := context.TODO()
	err := testDao.setCache(ctx, "test1", "VALUE", 10)
	assert.NoError(t, err)
	assert.Equal(t, testDao.getCacheString(ctx, "test1"), "VALUE")
}

func TestDao_DelCache(t *testing.T) {
	ctx := context.TODO()
	_ = testDao.setCache(ctx, "test1", 2>>100, 10)
	err := testDao.setCache(ctx, "test2", "2000", 10)
	assert.NoError(t, err)
	err = testDao.delCache(ctx, "test1", "test2")
	assert.NoError(t, err)
	assert.Equal(t, testDao.getCacheInt64(ctx, "test1"), int64(0))
	assert.Equal(t, testDao.getCacheString(ctx, "test2"), "")
}
