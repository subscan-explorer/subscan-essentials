package dao

const (
	TransferModule = "transfer"
)

var (
	RedisMetadataKey           = redisKeyPrefix() + "metadata"
	blockCacheKey              = redisKeyPrefix() + "block:%d"
	blockByHashCacheKey        = redisKeyPrefix() + "blockByHash:%s"
	RedisFillAlreadyBlockNum   = redisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum = redisKeyPrefix() + "FillFinalizedBlockNum"
	RedisRepairBlockKey        = redisKeyPrefix() + "RepairBlock"
)
