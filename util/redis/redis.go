package redis

import (
	"github.com/itering/subscan/configs"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	SubPool *redis.Pool
	once    sync.Once
)

func Init() {
	once.Do(func() {
		redisAddr := configs.Boot.Redis.Addr
		redisPassword := configs.Boot.Redis.Password
		redisDatabase := configs.Boot.Redis.DbName
		SubPool = initRedis(redisDatabase, redisPassword, redisAddr)
	})
}

func initRedis(database int, password string, server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialPassword(password), redis.DialDatabase(database))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
