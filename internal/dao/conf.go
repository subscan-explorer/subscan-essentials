package dao

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/itering/subscan/internal/util"
)

type (
	mysqlConf struct {
		Api  *sql.Config
		Task *sql.Config
		Test *sql.Config
	}
	redisConf struct {
		Config *redis.Config
		DbName int
	}
)

func (dc *mysqlConf) mergeEnvironment() {
	dbHost := util.GetEnv("MYSQL_HOST", "127.0.0.1")
	dbUser := util.GetEnv("MYSQL_USER", "root")
	dbPass := util.GetEnv("MYSQL_PASS", "")
	dbName := util.GetEnv("MYSQL_DB", util.NetworkNode)
	dc.Api.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName) + dc.Api.DSN
	dc.Task.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName) + dc.Task.DSN
}

func (rc *redisConf) mergeEnvironment() {
	rc.Config.Addr = util.GetEnv("REDIS_HOST", "127.0.0.1") + ":" + util.GetEnv("REDIS_PORT", "6379")
	rc.DbName = util.StringToInt(util.GetEnv("REDIS_DATABASE", "0"))
}
