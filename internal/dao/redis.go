package dao

import (
	"context"
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
