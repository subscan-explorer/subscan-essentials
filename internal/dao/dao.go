package dao

import (
	"context"
	"github.com/itering/substrate-api-rpc/websocket"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	"github.com/jinzhu/gorm"
)

var (
	DaemonAction = []string{"substrate"}
)

// dao
type Dao struct {
	db    *gorm.DB
	redis *redis.Pool
	cache *fanout.Fanout
}

func (d *Dao) SpecialMetadata(spec int) string {
	if metadata := d.RuntimeVersionRaw(spec); metadata != nil {
		return metadata.Raw
	}
	return ""
}

func (d *Dao) RPCPool() *websocket.PoolConn {
	conn, _ := websocket.Init()
	return conn
}

func (d *Dao) DB() *gorm.DB {
	return d.db
}

func (d *Dao) Redis() *redis.Pool {
	return d.redis
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
