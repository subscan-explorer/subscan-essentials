package dao

import (
	"context"
	"time"

	"log"
)

func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn, _ := d.redis.Redis().GetContext(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Printf("conn.Set(PING) error(%v)", err)
	}
	return
}

func (d *Dao) SetHeartBeatNow(c context.Context, action string) error {
	return d.redis.SetCache(c, action, time.Now().Unix(), 300)
}

// private funcs
