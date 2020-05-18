package dao

const (
	TransferModule   = "transfer"
	RingBlanceModule = "Balances"
	Kton             = "kton"
)

var (
	enableBalancesModule     = map[string]string{"balances": RingBlanceModule, "kton": Kton}
	RedisMetadataKey         = redisKeyPrefix() + "metadata"
	blockCacheKey            = redisKeyPrefix() + "block:%d"
	blockByHashCacheKey      = redisKeyPrefix() + "blockByHash:%s"
	RedisFillAlreadyBlockNum = redisKeyPrefix() + "FillAlreadyBlockNum"
	RedisRepairBlockKey      = redisKeyPrefix() + "RepairBlock"
)
