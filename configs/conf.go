package configs

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/itering/subscan/util"
)

type (
	MysqlConf struct {
		Conf struct {
			Host string
			User string
			Pass string
			DB   string
		}
		Api  *sql.Config
		Task *sql.Config
		Test *sql.Config
	}
	RedisConf struct {
		Config *redis.Config
		DbName int
	}
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (dc *MysqlConf) MergeConf() {
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(dc))
	dc.mergeEnvironment()
}

func (rc *RedisConf) MergeConf() {
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(rc))
	rc.mergeEnvironment()
}

func (dc *MysqlConf) mergeEnvironment() {
	dbHost := util.GetEnv("MYSQL_HOST", dc.Conf.Host)
	dbUser := util.GetEnv("MYSQL_USER", dc.Conf.User)
	dbPass := util.GetEnv("MYSQL_PASS", dc.Conf.Pass)
	dbName := util.GetEnv("MYSQL_DB", dc.Conf.DB)
	dc.Api.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName) + dc.Api.DSN
	dc.Task.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName) + dc.Task.DSN
}

func (rc *RedisConf) mergeEnvironment() {
	rc.Config.Addr = util.GetEnv("REDIS_ADDR", rc.Config.Addr)
	rc.DbName = util.StringToInt(util.GetEnv("REDIS_DATABASE", "0"))
}
