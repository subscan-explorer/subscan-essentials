package dao

var (
	RedisMetadataKey           = redisKeyPrefix() + "metadata"
	RedisFillAlreadyBlockNum   = redisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum = redisKeyPrefix() + "FillFinalizedBlockNum"
	RedisRepairBlockKey        = redisKeyPrefix() + "RepairBlock"
)
