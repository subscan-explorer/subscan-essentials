package redisDao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/configs"
	"time"
)

type Dao struct {
	redis *redis.Pool
}

func (d *Dao) Redis() *redis.Pool {
	return d.redis
}

func Init() *Dao {
	pool := newCachePool(configs.Boot.Redis.Addr, "")
	return &Dao{redis: pool}
}

func newCachePool(host, password string) *redis.Pool {
	var pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			// the redis protocol should probably be made sett-able
			c, err := redis.Dial("tcp", host, redis.DialReadTimeout(time.Millisecond*200), redis.DialConnectTimeout(time.Millisecond*200), redis.DialWriteTimeout(time.Millisecond*200))
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					_ = c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}
	return pool
}

func (d *Dao) Close() error {
	return d.redis.Close()
}
