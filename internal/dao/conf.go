package dao

import (
	"fmt"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/database/sql"
	xtime "github.com/bilibili/kratos/pkg/time"
	"subscan-end/utiles"
)

type (
	MysqlConf struct {
		Dev *sql.Config
	}
	RedisConf struct {
		Dev       *redis.Config
		DevExpire xtime.Duration
	}
)

func (dc *MysqlConf) mergeEnvironment() {
	dbHost := utiles.GetEnv("MYSQL_HOST", "127.0.0.1")
	dbUser := utiles.GetEnv("MYSQL_USER", "root")
	dbPass := utiles.GetEnv("MYSQL_PASS", "")
	dbName := utiles.GetEnv("MYSQL_DB", "subscan-end")
	dc.Dev.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName) + dc.Dev.DSN
}

func (rc *RedisConf) mergeEnvironment() {
	rc.Dev.Addr = utiles.GetEnv("REDIS_ADDR", rc.Dev.Addr)
}
