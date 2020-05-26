package dao

const (
	TransferModule = "transfer"
	BalanceModule  = "balances"
	Kton           = "kton"
	Ring           = "ring"
)

var (
	enableBalancesModule       = map[string]string{BalanceModule: BalanceModule, Kton: Kton, Ring: BalanceModule}
	RedisMetadataKey           = redisKeyPrefix() + "metadata"
	blockCacheKey              = redisKeyPrefix() + "block:%d"
	blockByHashCacheKey        = redisKeyPrefix() + "blockByHash:%s"
	RedisFillAlreadyBlockNum   = redisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum = redisKeyPrefix() + "FillFinalizedBlockNum"
	RedisRepairBlockKey        = redisKeyPrefix() + "RepairBlock"
)
