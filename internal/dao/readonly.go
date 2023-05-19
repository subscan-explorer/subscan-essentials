package dao

import (
	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/configs"
	"gorm.io/gorm"
)

type ReadOnlyDao struct {
	db    *gorm.DB
	redis *redis.Pool
}

func NewReadOnly() (IReadOnlyDao, *DbStorage) {
	db := newDb()

	pool := newCachePool(configs.Boot.Redis.Addr, "")
	dao := &ReadOnlyDao{
		db:    db,
		redis: pool,
	}
	return dao, &DbStorage{db: db, dao: dao}
}

func readOnlyWithDb(db *gorm.DB, redis *redis.Pool) ReadOnlyDao {
	return ReadOnlyDao{db: db, redis: redis}
}
