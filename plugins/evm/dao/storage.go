package dao

import (
	"context"
	"fmt"
	subscan_plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Storage struct {
	dao   storage.DB
	db    *gorm.DB
	redis subscan_plugin.RedisPool
}

var sg *Storage

func Init(db storage.DB, redis subscan_plugin.RedisPool) *Storage {
	if sg == nil {
		sg = &Storage{
			dao:   db,
			db:    db.GetDbInstance().(*gorm.DB),
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
		fmt.Println(s.db == nil)
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
