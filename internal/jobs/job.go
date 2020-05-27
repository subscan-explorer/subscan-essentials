package jobs

import (
	"github.com/freehere107/go-workers"
	"github.com/itering/subscan/internal/util"
)

// Init worker instant
// worker use redis connect
// namespace is NETWORK_NODE env
func Init() {
	addr := util.GetEnv("REDIS_HOST", "127.0.0.1") + ":" + util.GetEnv("REDIS_PORT", "6379")
	workers.Configure(map[string]string{
		"server":    addr,
		"database":  "0",
		"pool":      "30",
		"process":   "1",
		"namespace": util.NetworkNode,
	})
}
