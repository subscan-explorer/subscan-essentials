package dao

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"strconv"
)

func (d *Dao) setCache(c context.Context, key string, value interface{}, ttl int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var val string
	switch v := value.(type) {
	case string:
		val = v
	case int64:
		val = strconv.FormatInt(v, 10)
	case int:
		val = strconv.Itoa(v)
	default:
		b, _ := json.Marshal(v)
		if val = string(b); val == "null" {
			return
		}
	}
	if ttl <= 0 {
		_, err = conn.Do("set", key, val)
		return
	}
	_, err = conn.Do("setex", key, ttl, val)
	return
}

func (d *Dao) getCacheBytes(c context.Context, key string) []byte {
	conn := d.redis.Get(c)
	defer conn.Close()
	if cache, err := redis.Bytes(conn.Do("get", key)); err == nil {
		return cache
	}
	return nil
}

func (d *Dao) getCacheString(c context.Context, key string) string {
	conn := d.redis.Get(c)
	defer conn.Close()
	if cache, err := redis.String(conn.Do("get", key)); err == nil {
		return cache
	}
	return ""
}

func (d *Dao) getCacheInt64(c context.Context, key string) int64 {
	conn := d.redis.Get(c)
	defer conn.Close()
	if cache, err := redis.Int64(conn.Do("get", key)); err == nil {
		return cache
	}
	return 0
}

func (d *Dao) delCache(c context.Context, key ...string) error {
	if len(key) == 0 {
		return nil
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}
	for _, v := range key {
		args = args.Add(v)
	}
	_, err := conn.Do("del", args...)
	return err
}
