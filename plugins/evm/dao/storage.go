package dao

import (
	"context"
	subscan_plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan/util/mq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Storage struct {
	db    *gorm.DB
	redis subscan_plugin.RedisPool
}

var sg *Storage

func Init(db *gorm.DB, redis subscan_plugin.RedisPool) *Storage {
	if sg == nil {
		sg = &Storage{
			db:    db,
			redis: redis,
		}
	}
	return sg
}

func (s *Storage) AddOrUpdateItem(c context.Context, item interface{}, keys []string, updates ...string) *gorm.DB {
	var keyFields []clause.Column
	for _, key := range keys {
		keyFields = append(keyFields, clause.Column{Name: key})
	}
	if len(updates) > 0 {
		return s.db.WithContext(c).Clauses(clause.OnConflict{
			Columns:   keyFields,
			DoUpdates: clause.AssignmentColumns(updates),
		}).Create(item)
	}
	return s.db.WithContext(c).Clauses(clause.OnConflict{
		Columns:   keyFields,
		UpdateAll: true,
	}).Create(item)
}

func (s *Storage) Tables() []interface{} {
	return []interface{}{
		&Transaction{},
		&TransactionReceipt{},
		&Contract{},
		&Token{},
		&TokenHolder{},
		&TokensTransfers{},
		&EvmBlock{},
		&Erc721Holders{},
		&AbiMapping{},
		&ERC1155Item{},
		&ERC1155Holder{},
		&Account{},
	}

}

func Publish(queue, class string, args interface{}) error {
	if mq.Instant == nil {
		return nil
	}
	return mq.Instant.Publish(queue, class, args)
}
