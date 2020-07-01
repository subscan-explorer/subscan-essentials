package dao

import (
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/jinzhu/gorm"
)

var TestDao *Dao

func init() {
	if client, err := paladin.NewFile("../../configs"); err != nil {
		panic(err)
	} else {
		paladin.DefaultClient = client
	}
	var (
		dc mysqlConf
		rc redisConf
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))

	db, err := gorm.Open("mysql", dc.Test.DSN)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	const testRedisDb = 1

	TestDao = &Dao{
		db:    db,
		redis: redis.NewPool(rc.Config, redis.DialDatabase(testRedisDb)),
	}
	var tables []string
	err = db.Raw("show tables;").Pluck("Tables_in_subscan_test", &tables).Error
	if err != nil {
		panic(err)
	}
	// for _, value := range tables {
	// 	db.DropTable(value)
	// }
	TestDao.Migration()
}

func NewTestDao() {

}
