package redisDao

import (
	"context"
	"encoding/json"
	"github.com/itering/subscan/util"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

func (d *Dao) SetCache(c context.Context, key string, value interface{}, ttl int) (err error) {
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

func (d *Dao) GetCacheString(c context.Context, key string) string {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	if cache, err := redis.String(conn.Do("get", key)); err == nil {
		return cache
	}
	return ""
}

func (d *Dao) GetCacheInt64(c context.Context, key string) int64 {
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

func (d *Dao) DelCache(c context.Context, key ...string) error {
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
	_ = d.SetCache(c, key, value, ttl)
	return nil
}

func (d *Dao) HMGet(c context.Context, key string, field ...string) (ms map[string]string) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	args := redis.Args{}
	args = args.Add(key)
	for _, v := range field {
		args = args.Add(v)
	}
	args = args.Add(c)
	rsp, _ := redis.Strings(conn.Do("HMGET", args...))
	if len(rsp) == 0 || len(rsp) != len(field) {
		return
	}
	ms = make(map[string]string)
	for i, k := range field {
		if v := rsp[i]; v != "" {
			ms[k] = rsp[i]
		}
	}
	return
}

func (d *Dao) GetCacheTtl(c context.Context, key string) int {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	if cache, err := redis.Int(conn.Do("ttl", key, c)); err == nil {
		return cache
	}
	return -2
}

func (d *Dao) HmSet(c context.Context, key string, value interface{}) (err error) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	args := redis.Args{}.Add(key)
	switch v := value.(type) {
	case map[string]string:
		for k, v := range v {
			args = args.Add(k).Add(v)
		}
	case map[string]int:
		for k, v := range v {
			args = args.Add(k).Add(v)
		}
	case map[string][]string:
		for k, v := range v {
			args = args.Add(k).Add(util.ToString(v))
		}
	}
	if len(args) <= 1 {
		return
	}
	// append context to the args end
	args.Add(c)
	_, err = conn.Do("HMSET", args...)
	return
}

func (d *Dao) HmSetEx(c context.Context, key string, value interface{}, ttl int) (err error) {
	conn, _ := d.redis.GetContext(c)
	defer conn.Close()
	defer func() {
		if ttl > 0 {
			_ = conn.Send("EXPIRE", key, ttl)
		}
	}()
	err = d.HmSet(c, key, value)
	return err
}
