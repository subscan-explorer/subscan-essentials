package jobs

import (
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/freehere107/go-workers"
	"subscan-end/internal/dao"
)

func Init() {
	var rc dao.RedisConf
	_ = paladin.Get("redis.toml").UnmarshalTOML(&rc)
	workers.Configure(map[string]string{
		"server":   rc.Dev.Addr,
		"database": "0",
		"pool":     "30",
		"process":  "1",
	})
}
