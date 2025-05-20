package dao

import "github.com/itering/subscan/model"

var (
	RedisMetadataKey           = model.RedisKeyPrefix() + "metadata"
	RedisFillAlreadyBlockNum   = model.RedisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum = model.RedisKeyPrefix() + "FillFinalizedBlockNum"
)

// local cache value
var maxTableBlockNum uint = 0
