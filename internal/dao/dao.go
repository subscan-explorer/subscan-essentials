package dao

import (
	"context"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var (
	DaemonAction = []string{"substrate", "worker", "cronWorker"}
	Provider     = wire.NewSet(New)
)

// Dao dao.
type Dao struct {
	db    *gorm.DB
	redis *redis.Pool
	cache *fanout.Fanout
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// New new a dao and return.
func New() (dao *Dao) {
	var (
		dc mysqlConf
		rc redisConf
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	dc.mergeEnvironment()
	rc.mergeEnvironment()
	dao = &Dao{
		db:    newDb(dc),
		redis: redis.NewPool(rc.Config, redis.DialDatabase(rc.DbName)),
		cache: fanout.New("scan", fanout.Worker(1), fanout.Buffer(1024)),
	}
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
