package redisDao

import (
	"context"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RedisCache(t *testing.T) {
	ctx := context.TODO()
	util.ConfDir = "../../configs"
	configs.Init()
	dao := Init()
	t.Run("SetCache", func(t *testing.T) {
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
			err := dao.SetCache(ctx, test.key, test.value, 10)
			assert.NoError(t, err)
			assert.Equal(t, dao.GetCacheString(ctx, test.key), test.expect)
		}
	})

	t.Run("GetCacheBytes", func(t *testing.T) {
		err := dao.SetCache(ctx, "test1", "value", 10)
		assert.NoError(t, err)
		assert.Equal(t, dao.GetCacheBytes(ctx, "test1"), []byte{0x76, 0x61, 0x6c, 0x75, 0x65})
	})

	t.Run("GetCacheInt64", func(t *testing.T) {
		err := dao.SetCache(ctx, "test1", 2000, 10)
		assert.NoError(t, err)
		assert.Equal(t, dao.GetCacheInt64(ctx, "test1"), int64(2000))
	})

	t.Run("DelCache", func(t *testing.T) {
		_ = dao.SetCache(ctx, "test1", 2>>100, 10)
		err := dao.SetCache(ctx, "test2", "2000", 10)
		assert.NoError(t, err)
		err = dao.DelCache(ctx, "test1", "test2")
		assert.NoError(t, err)
		assert.Equal(t, dao.GetCacheInt64(ctx, "test1"), int64(0))
		assert.Equal(t, dao.GetCacheString(ctx, "test2"), "")
	})

	t.Run("GetCacheString", func(t *testing.T) {
		err := dao.SetCache(ctx, "test1", "VALUE", 10)
		assert.NoError(t, err)
		assert.Equal(t, dao.GetCacheString(ctx, "test1"), "VALUE")
	})

}
