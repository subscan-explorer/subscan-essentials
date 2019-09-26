package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bilibili/kratos/pkg/sync/pipeline/fanout"
	"github.com/jinzhu/gorm"
	"strings"
	"subscan-end/internal/model"
	"subscan-end/libs/substrate"
	"subscan-end/utiles"
	"subscan-end/utiles/es"
	"time"

	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
)

var DaemonAction = []string{"substrate", "worker"}

// Dao dao.
type Dao struct {
	db           *gorm.DB
	redis        *redis.Pool
	redisExpire  int32
	substrateApi *substrate.Websocket
	es           *es.EsClient
	cache        *fanout.Fanout
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
		dc MysqlConf
		rc RedisConf
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	dc.mergeEnvironment()
	rc.mergeEnvironment()
	esClient, _ := es.NewEsClient()
	dao = &Dao{
		// mysql
		db: initDb(dc),
		// redis
		redis:       redis.NewPool(rc.Dev),
		redisExpire: int32(time.Duration(rc.DevExpire) / time.Second),
		cache:       fanout.New("scan", fanout.Worker(1), fanout.Buffer(1024)),
		// substrate rpc
		substrateApi: &substrate.Websocket{Provider: utiles.ProviderEndPoint},
		// es
		es: esClient,
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	_ = d.redis.Close()
	_ = d.db.Close()
	if d.es != nil {
		d.es.Client.Stop()
	}
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	//gorm auto ping
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

func (d *Dao) Migration() {
	db := d.db
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.ChainBlock{},
		&model.ChainEvent{},
		&model.ChainExtrinsic{},
		&model.ChainTransaction{},
		&model.ChainAccount{},
		&model.RuntimeVersion{},
		&model.ChainLog{},
		&model.DailyStatic{},
		&model.ChainSession{},
		&model.SessionValidator{},
		&model.SessionNominator{},
		&model.ValidatorInfo{},
	)
	db.Model(model.ChainTransaction{}).AddIndex("from_hex", "from_hex")
	db.Model(model.ChainTransaction{}).AddIndex("destination", "destination")
	db.Model(model.ChainTransaction{}).AddIndex("call_module_function", "call_module_function")
	db.Model(model.ChainBlock{}).AddUniqueIndex("hash", "hash")
	db.Model(model.ChainBlock{}).AddUniqueIndex("block_num", "block_num")
	db.Model(model.ChainBlock{}).AddIndex("codec_error", "codec_error")
	db.Model(model.ChainAccount{}).AddUniqueIndex("address", "address")
	db.Model(model.ChainExtrinsic{}).AddIndex("extrinsic_hash", "extrinsic_hash")
	db.Model(model.ChainExtrinsic{}).AddUniqueIndex("extrinsic_index", "extrinsic_index")
	db.Model(model.ChainExtrinsic{}).AddIndex("block_num", "block_num")
	db.Model(model.ChainExtrinsic{}).AddIndex("is_signed", "is_signed")
	db.Model(model.ChainExtrinsic{}).AddIndex("account_id", "is_signed,account_id")
	db.Model(model.ChainEvent{}).AddIndex("block_num", "block_num")
	db.Model(model.ChainEvent{}).AddIndex("type", "type")
	db.Model(model.ChainEvent{}).AddIndex("event_index", "event_index")
	db.Model(model.ChainLog{}).AddUniqueIndex("log_index", "log_index")
	db.Model(model.ChainLog{}).AddIndex("block_num", "block_num")
	db.Model(model.RuntimeVersion{}).AddUniqueIndex("spec_version", "spec_version")
	db.Model(model.DailyStatic{}).AddUniqueIndex("time_utc", "time_utc")
	db.Model(model.ChainSession{}).AddUniqueIndex("session_id", "session_id")
	db.Model(model.SessionValidator{}).AddUniqueIndex("session_rank", "session_id", "rank_validator")
	db.Model(model.SessionNominator{}).AddUniqueIndex("session_rank", "session_id", "rank_validator", "rank_nominator")
	db.Model(model.ValidatorInfo{}).AddUniqueIndex("validator_controller", "validator_controller")
	db.Model(model.ValidatorInfo{}).AddUniqueIndex("validator_stash", "validator_stash")
}

func initDb(dc MysqlConf) *gorm.DB {
	db, err := gorm.Open("mysql", dc.Dev.DSN)
	if err != nil {
		panic(err)
	}
	db.DB().SetConnMaxLifetime(5 * time.Minute)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	db.SetLogger(ormLog{})
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
		cacheKey := fmt.Sprintf("%s:heartBeat:%s", redisKeyPrefix(), dt)
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

func (d *Dao) BroadCastToChanel(c context.Context, topic string, msg interface{}) {
	conn := d.redis.Get(c)
	defer conn.Close()
	bs, _ := json.Marshal(map[string]interface{}{"topic": topic, "content": msg})
	if _, err := conn.Do("publish", utiles.SubScanChannel, string(bs)); err != nil {
		log.Error(err.Error())
	}
}

func redisKeyPrefix() string {
	return utiles.NetworkNode + ":"
}
