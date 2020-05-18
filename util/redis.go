package util

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var (
	SubPool        *redis.Pool
	SubScanChannel = fmt.Sprintf("%s_scan", NetworkNode)

	// subscribe endpoint
	wsEndpointCacheKey = NetworkNode + ":" + "wsEndpoint"
	PubTopicMetadata   = "metadata_update"
	PubNewFinalized    = "new_finalized"
)

func init() {
	redisAddr := GetEnv("REDIS_HOST", "127.0.0.1")
	redisPort := GetEnv("REDIS_PORT", "6379")
	redisPassword := GetEnv("REDIS_PASSWORD", "")
	redisDatabase := GetEnv("REDIS_DATABASE", "0")
	SubPool = initRedis(redisDatabase, redisPassword, redisAddr+":"+redisPort)
}

func initRedis(database string, password string, server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			db, _ := strconv.Atoi(database)
			c, err := redis.Dial("tcp", server, redis.DialPassword(password), redis.DialDatabase(db))
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

func SetEndpointCache(endpoint string) {
	conn := SubPool.Get()
	defer conn.Close()
	_, _ = conn.Do("set", wsEndpointCacheKey, endpoint)
}

func WsEndpointCache() string {
	conn := SubPool.Get()
	defer conn.Close()
	endpoint, _ := redis.String(conn.Do("get", wsEndpointCacheKey))
	if endpoint == "" {
		return WSEndPoint
	}
	return endpoint
}
