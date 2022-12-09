package dao

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/configs"
	"github.com/jinzhu/gorm"
)

var (
	DaemonAction = []string{"substrate"}
)

// dao
type Dao struct {
	db    *gorm.DB
	redis *redis.Pool
}

// New new a dao and return.
func New() (dao *Dao, storage *DbStorage) {
	db := newDb()

	dao = &Dao{
		db:    db,
		redis: newCachePool(configs.Boot.Redis.Addr, ""),
	}
	dao.Migration()
	storage = &DbStorage{db: db}
	return
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

// Close close the resource.
func (d *Dao) Close() {
	if d.redis != nil {
		_ = d.redis.Close()
	}
	_ = d.db.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	// gorm auto ping
	return
}
