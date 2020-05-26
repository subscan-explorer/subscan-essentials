package dao

import "fmt"

const (
	TransferModule = "transfer"
	BalanceModule  = "balances"
	Kton           = "kton"
	Ring           = "ring"

	validatorRole = "validator"
	nominatorRole = "nominator"
)

var (
	enableBalancesModule          = map[string]string{BalanceModule: BalanceModule, Kton: Kton, Ring: BalanceModule}
	symbolModuleMap               = map[string]string{"KTON": Kton}
	RedisMetadataKey              = redisKeyPrefix() + "metadata"
	blockCacheKey                 = redisKeyPrefix() + "block:%d"
	blockByHashCacheKey           = redisKeyPrefix() + "blockByHash:%s"
	RedisFillAlreadyBlockNum      = redisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum    = redisKeyPrefix() + "FillFinalizedBlockNum"
	RedisRepairBlockKey           = redisKeyPrefix() + "RepairBlock"
	ktonPool                      = redisKeyPrefix() + "KtonPool"
	RingPool                      = redisKeyPrefix() + "RingPool"
	totalIssuanceKey              = redisKeyPrefix() + "totalIssuance:%s"
	validatorCountKey             = redisKeyPrefix() + "validatorCount"
	validatorAddressKey           = redisKeyPrefix() + "validatorAddress"
	validatorControllerAddressKey = redisKeyPrefix() + "validatorControllerAddress"
	nominatorAddressKey           = redisKeyPrefix() + "nominatorAddress"
	eraStartBlockKey              = redisKeyPrefix() + "eraStartBlock"
)

func (d *Dao) CacheGetTotalIssuanceKey(module string) string {
	return fmt.Sprintf(totalIssuanceKey, module)
}
