package dao

import (
	"context"
	"fmt"
	redisDao "github.com/itering/subscan/share/redis"
	"github.com/itering/subscan/util"
	"gorm.io/gorm"
	"os"
	"time"
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
	if os.Getenv("MYSQL_MIGRATE") != "false" {
		util.Logger().Info("start db auto migrate")
		now := time.Now()
		dao.Migration()
		util.Logger().Info(fmt.Sprintf("db auto migrate complete %s %d ms", "cost", time.Since(now).Milliseconds()))
	}
	storage = &DbStorage{db: db, DbDriver: dao.DbDriver, d: dao}
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
