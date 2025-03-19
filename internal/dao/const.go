package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	RedisMetadataKey           = redisKeyPrefix() + "metadata"
	RedisFillAlreadyBlockNum   = redisKeyPrefix() + "FillAlreadyBlockNum"
	RedisFillFinalizedBlockNum = redisKeyPrefix() + "FillFinalizedBlockNum"
)

// local cache value
var maxTableBlockNum uint = 0

func IgnoreDuplicate(tx *gorm.DB) *gorm.DB {
	return tx.Clauses(clause.OnConflict{DoNothing: true})
}
