package dao

var (
	RedisMetadataKey           = redisKeyPrefix() + "metadata"
	RedisFillAlreadyBlockNum   = redisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum = redisKeyPrefix() + "FillFinalizedBlockNum"
)

// local cache value
var maxTableBlockNum uint = 0
