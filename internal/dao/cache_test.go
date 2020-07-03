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
		err := TestDao.SetCache(ctx, test.key, test.value, 10)
		assert.NoError(t, err)
		assert.Equal(t, TestDao.GetCacheString(ctx, test.key), test.expect)
	}
}

func TestDao_GetCacheBytes(t *testing.T) {
	ctx := context.TODO()
	err := TestDao.SetCache(ctx, "test1", "value", 10)
	assert.NoError(t, err)
	assert.Equal(t, TestDao.GetCacheBytes(ctx, "test1"), []byte{0x76, 0x61, 0x6c, 0x75, 0x65})
}

func TestDao_GetCacheInt64(t *testing.T) {
	ctx := context.TODO()
	err := TestDao.SetCache(ctx, "test1", 2000, 10)
	assert.NoError(t, err)
	assert.Equal(t, TestDao.GetCacheInt64(ctx, "test1"), int64(2000))
}

func TestDao_GetCacheString(t *testing.T) {
	ctx := context.TODO()
	err := TestDao.SetCache(ctx, "test1", "VALUE", 10)
	assert.NoError(t, err)
	assert.Equal(t, TestDao.GetCacheString(ctx, "test1"), "VALUE")
}

func TestDao_DelCache(t *testing.T) {
	ctx := context.TODO()
	_ = TestDao.SetCache(ctx, "test1", 2>>100, 10)
	err := TestDao.SetCache(ctx, "test2", "2000", 10)
	assert.NoError(t, err)
	err = TestDao.DelCache(ctx, "test1", "test2")
	assert.NoError(t, err)
	assert.Equal(t, TestDao.GetCacheInt64(ctx, "test1"), int64(0))
	assert.Equal(t, TestDao.GetCacheString(ctx, "test2"), "")
}
