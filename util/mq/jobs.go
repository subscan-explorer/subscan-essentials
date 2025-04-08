package mq

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/itering/subscan/util"
	redisUtil "github.com/itering/subscan/util/redis"

	"github.com/gomodule/redigo/redis"
)

type IJob interface {
	Type() string
	Init()
	Consumption()
	Publish(string, string, interface{}) error
	ForcePublish(string, string, interface{}) error
	Shutdown(_ context.Context) error
	SubscribePublish(any) error
}

func rateLimit(c context.Context, queue, class string, args interface{}) bool {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", queue, class, util.ToString(args))))
	formatKey := fmt.Sprintf("%s:rateLimit:%x", util.NetworkNode, hash)
	conn, _ := redisUtil.SubPool.GetContext(c)
	defer conn.Close() // nolint: errcheck
	if ttl, _ := redis.Int64(conn.Do("ttl", formatKey)); ttl > 0 {
		return true
	}
	_, _ = conn.Do("setex", formatKey, 12, "1")
	return false
}

var (
	Instant IJob
)

const (
	GoWorkerName = "go-worker"
)

type CommonMessage struct {
	Data  []byte `json:"data"`
	Queue string `json:"queue"`
	Class string `json:"class"`
}

func New() IJob {
	Instant = &GoWorker{}
	Instant.Init()
	return Instant
}
