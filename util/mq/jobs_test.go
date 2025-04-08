package mq

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/util"
	redisUtil "github.com/itering/subscan/util/redis"
	"testing"
)

func Test_ratelimit(t *testing.T) {
	util.ConfDir = "../../configs"
	configs.Init()
	redisUtil.Init()

	t.Run("Test limit cache key will set ", func(t *testing.T) {
		c := context.Background()
		queue := "testQueue"
		class := "testClass"
		args := "testArgs"

		// Set up a mock Redis connection
		conn := redisUtil.SubPool.Get()
		defer conn.Close()

		// Clear the key before testing
		hash := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", queue, class, util.ToString(args))))
		formatKey := fmt.Sprintf("%s:rateLimit:%x", util.NetworkNode, hash)
		_, _ = conn.Do("DEL", formatKey)

		// First call should return false
		if rateLimit(c, queue, class, args) {
			t.Errorf("Expected rateLimit to return false on first call")
		}

		// Second call should return true
		if !rateLimit(c, queue, class, args) {
			t.Errorf("Expected rateLimit to return true on second call")
		}
	})
}
