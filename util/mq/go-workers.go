package mq

import (
	"context"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/util"

	"github.com/itering/go-workers"
)

type GoWorker struct{}

func (g *GoWorker) SubscribePublish(any) error {
	return nil
}

func (g *GoWorker) Type() string {
	return GoWorkerName
}

func (g *GoWorker) Publish(queue, class string, args interface{}) error {
	if rateLimit(context.TODO(), queue, class, args) {
		return nil
	}
	return g.ForcePublish(queue, class, args)
}

// Init worker instant
// worker use redis connect
// namespace is NETWORK_NODE env
func (g *GoWorker) Init() {
	workers.Configure(map[string]string{
		"server":    configs.Boot.Redis.Addr,
		"database":  util.IntToString(configs.Boot.Redis.DbName),
		"pool":      "30",
		"process":   "1",
		"namespace": util.NetworkNode,
	})
}

func (g *GoWorker) Shutdown(context.Context) error {
	return nil
}

func (g *GoWorker) Consumption() {
	workers.Run()
}

func (g *GoWorker) ForcePublish(queue, class string, args interface{}) error {
	if _, err := workers.Enqueue(queue, class, args); err != nil {
		return err
	}
	return nil
}
