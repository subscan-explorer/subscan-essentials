package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/configs"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

var DaemonAction = []string{"substrate"}

// dao
type Dao struct {
	db    *gorm.DB
	redis *redis.Pool
	ReadOnlyDao
}

// New new a dao and return.
func New(migrate bool) (IDao, *DbStorage) {
	db := newDb()

	pool := newCachePool(configs.Boot.Redis.Addr, "")
	dao := &Dao{
		db:          db,
		redis:       pool,
		ReadOnlyDao: readOnlyWithDb(db, pool),
	}
	dao.Migration()
	storage := &DbStorage{db: db, dao: dao}
	return dao, storage
}

func newCachePool(host, password string) *redis.Pool {
	pool := &redis.Pool{
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
	slog.Debug("redis pool init success")
	return pool
}

// Close close the resource.
func (d *ReadOnlyDao) Close() {
	if d.redis != nil {
		_ = d.redis.Close()
	}
	db, _ := d.db.DB()
	_ = db.Close()
}

// Ping ping the resource.
func (d *ReadOnlyDao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	// gorm auto ping
	return
}

func (d *ReadOnlyDao) GetModelTableName(model interface{}) string {
	stmt := &gorm.Statement{DB: d.db}
	if err := stmt.Parse(model); err != nil {
		panic(fmt.Sprintf("get model table name error: %v", err))
	}
	return stmt.Schema.Table
}

func where(query string, args ...interface{}) whereClauses {
	return whereClauses{query: query, args: args}
}

type whereClauses struct {
	query string
	args  []interface{}
}

func findOne[T any](d *ReadOnlyDao, sel string, where whereClauses, orderBy interface{}) (*T, error) {
	var find []T
	if sel == "" {
		sel = "*"
	}

	tx := d.db.Select(sel).Where(where.query, where.args...).Limit(1)
	if orderBy != nil {
		tx = tx.Order(orderBy)
	}
	res := tx.Find(&find)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(find) == 0 {
		return nil, nil
	}
	return &find[0], nil
}
