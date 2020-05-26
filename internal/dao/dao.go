package dao

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	"github.com/go-sql-driver/mysql"
	"github.com/itering/subscan/util"
	"github.com/jinzhu/gorm"
)

var DaemonAction = []string{"substrate", "worker", "cronWorker"}

// Dao dao.
type Dao struct {
	db          *gorm.DB
	redis       *redis.Pool
	redisExpire int32
	cache       *fanout.Fanout
}

type GormDB struct {
	*gorm.DB
	gdbDone bool
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(strings.Repeat("%v ", len(v)), v...)
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
	// esClient, _ := es.NewEsClient()
	dao = &Dao{
		db:    initDb(dc),
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

func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

func initDb(dc mysqlConf) (db *gorm.DB) {
	var err error
	if os.Getenv("TASK_MOD") == "true" {
		db, err = gorm.Open("mysql", dc.Task.DSN)
	} else if os.Getenv("TEST_MOD") == "true" {
		db, err = gorm.Open("mysql", dc.Test.DSN)
	} else {
		db, err = gorm.Open("mysql", dc.Api.DSN)
	}
	if err != nil {
		panic(err)
	}
	db.DB().SetConnMaxLifetime(5 * time.Minute)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	if util.IsProduction() {
		db.SetLogger(ormLog{})
	}
	db.LogMode(true)
	return db
}

func (d *Dao) SetHeartBeatNow(c context.Context, action string) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, _ = conn.Do("SET", action, time.Now().Unix())
}

func (d *Dao) GetHeartBeatNow(c context.Context) map[string]bool {
	conn := d.redis.Get(c)
	defer conn.Close()
	status := map[string]bool{}
	for _, dt := range DaemonAction {
		cacheKey := fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, dt)
		t, err := redis.Int64(conn.Do("get", cacheKey))
		if err != nil || time.Now().Unix()-t > 60 {
			status[dt] = false
		} else {
			status[dt] = true
		}
	}
	return status
}

func (d *Dao) DbBegin() *GormDB {
	txn := d.db.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	return &GormDB{txn, false}
}

func (c *GormDB) DbCommit() {
	if c.gdbDone {
		return
	}
	tx := c.Commit()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		fmt.Println("Fatal error DbCommit", err)
	}
}

func (c *GormDB) DbRollback() {
	if c.gdbDone {
		return
	}
	tx := c.Rollback()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		fmt.Println("Fatal error DbRollback", err)
	}
}

func redisKeyPrefix() string {
	return util.NetworkNode + ":"
}

func (d *Dao) checkDBError(err error) error {
	if err == mysql.ErrInvalidConn || err == driver.ErrBadConn {
		return err
	}
	return nil
}
