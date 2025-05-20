package dao

import (
	"context"
	redisDao "github.com/itering/subscan/share/redis"
	"gorm.io/gorm"
)

type Dao struct {
	db       *gorm.DB
	redis    *redisDao.Dao
	DbDriver string
}

// New new a dao and return.
func New() (dao *Dao, storage *DbStorage, pool *redisDao.Dao) {
	db := newDb()
	pool = redisDao.Init()
	dao = &Dao{
		db:       db,
		redis:    pool,
		DbDriver: db.Dialector.Name(),
	}
	dao.Migration()
	storage = &DbStorage{db: db, DbDriver: dao.DbDriver}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	if d.redis != nil {
		_ = d.redis.Close()
	}
	_ = d.db
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	// gorm auto ping
	return
}
