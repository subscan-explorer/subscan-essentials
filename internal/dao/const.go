package dao

import "github.com/itering/subscan/model"

var (
	RedisMetadataKey           = model.RedisKeyPrefix() + "metadata"
	RedisFillAlreadyBlockNum   = model.RedisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum = model.RedisKeyPrefix() + "FillFinalizedBlockNum"
	RedisExtrinsicCountKey     = model.RedisKeyPrefix() + "extrinsic_count"
)

// local cache value
var maxTableBlockNum uint = 0
