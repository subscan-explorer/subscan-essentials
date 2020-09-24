package dao

import (
	"context"
	"github.com/go-kratos/kratos/pkg/cache/redis"
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
	var dc configs.MysqlConf
	var rc configs.RedisConf
	dc.MergeConf()
	rc.MergeConf()
	db := newDb(dc)
	dao = &Dao{
		db:    db,
		redis: redis.NewPool(rc.Config, redis.DialDatabase(rc.DbName)),
	}
	dao.Migration()
	storage = &DbStorage{db: db}
	return
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
