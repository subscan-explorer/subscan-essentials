package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/itering/subscan/util"
	"golang.org/x/exp/slog"
)

func (d *ReadOnlyDao) pingRedis(ctx context.Context) (err error) {
	conn, _ := d.redis.GetContext(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		slog.Debug("conn.Set(PING) error(%v)", err)
	}
	return
}

func (d *Dao) SetHeartBeatNow(c context.Context, action string) error {
	return d.setCache(c, action, time.Now().Unix(), 300)
}

func (d *ReadOnlyDao) DaemonHealth(c context.Context) map[string]bool {
	conn, _ := d.redis.GetContext(c)
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

// private funcs
func redisKeyPrefix() string {
	return util.NetworkNode + ":"
}
