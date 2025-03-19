package dao

import (
	"context"
	"encoding/json"
	"github.com/itering/subscan/util"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

func (d *Dao) setCache(c context.Context, key string, value interface{}, ttl int) (err error) {
	conn, _ := d.redis.GetContext(c)
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
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	if cache, err := redis.Bytes(conn.Do("get", key)); err == nil {
		return cache
	}
	return nil
}

func (d *Dao) getCacheString(c context.Context, key string) string {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	if cache, err := redis.String(conn.Do("get", key)); err == nil {
		return cache
	}
	return ""
}

func (d *Dao) getCacheInt64(c context.Context, key string) int64 {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	if cache, err := redis.Int64(conn.Do("get", key)); err == nil {
		return cache
	}
	return 0
}

func (d *Dao) GetCacheBytes(c context.Context, key string) []byte {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	if cache, err := redis.Bytes(conn.Do("get", key)); err == nil {
		return cache
	}
	return nil
}

func (d *Dao) delCache(c context.Context, key ...string) error {
	if len(key) == 0 {
		return nil
	}
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	args := redis.Args{}
	for _, v := range key {
		args = args.Add(v)
	}
	_, err := conn.Do("del", args...)
	return err
}

// FetchCache try fetch a cache, do action if cache not exist
// temp add opt params force
func (d *Dao) FetchCache(c context.Context, key string, value interface{}, action func() error, ttl int, force ...bool) error {
	if len(force) == 0 || !force[0] {
		if b := d.GetCacheBytes(c, key); b != nil {
			return util.UnmarshalAny(value, b)
		}
	}
	if err := action(); err != nil {
		return err
	}
	_ = d.setCache(c, key, value, ttl)
	return nil
}
