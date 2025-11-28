package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"strings"
)

func (d *Dao) CreateEvent(txn *GormDB, events []model.ChainEvent) error {
	if len(events) == 0 {
		return nil
	}
	query := txn.Scopes(d.TableNameFunc(events[0]), model.IgnoreDuplicate).CreateInBatches(events, 2000)
	return query.Error
}

// GetEventListCursor implements bidirectional cursor pagination on events using id as cursor.
func (d *Dao) GetEventListCursor(ctx context.Context, limit int, _ string, fixedTableIndex int, beforeId uint, afterId uint, where ...model.Option) (list []model.ChainEvent, hasPrev, hasNext bool) {
	fetchLimit := limit + 1
	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	maxTableIndex := blockNum / int(model.SplitTableBlockNum)
	if afterId > 0 {
		maxTableIndex = int(afterId/model.SplitTableBlockNum) / model.IdGenerateCoefficient
	}
	if fixedTableIndex >= 0 {
		maxTableIndex = fixedTableIndex
	}

	if afterId > 0 { // next page
		for index := maxTableIndex; index >= 0 && len(list) < fetchLimit; index-- {
			if fixedTableIndex >= 0 && index != fixedTableIndex {
				continue
			}
			var tableData []model.ChainEvent
			q := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainEvent{BlockNum: uint(index) * model.SplitTableBlockNum}))
			q = q.Scopes(where...).Where("id < ?", afterId).Order("id desc").Limit(fetchLimit - len(list))
			if err := q.Find(&tableData).Error; err != nil {
				continue
			}
			list = append(list, tableData...)
		}
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = true
		return
	}

	if beforeId > 0 { // previous page
		startIdx := int(beforeId/model.SplitTableBlockNum) / model.IdGenerateCoefficient
		if fixedTableIndex >= 0 {
			startIdx = fixedTableIndex
		}
		for index := startIdx; index <= maxTableIndex && len(list) < fetchLimit; index++ {
			if fixedTableIndex >= 0 && index != fixedTableIndex {
				continue
			}
			var tableData []model.ChainEvent
			q := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainEvent{BlockNum: uint(index) * model.SplitTableBlockNum}))
			q = q.Scopes(where...)
			if index == startIdx {
				q = q.Where("id > ?", beforeId)
			}
			q = q.Order("id asc").Limit(fetchLimit - len(list))
			if err := q.Find(&tableData).Error; err != nil {
				continue
			}
			list = append(list, tableData...)
		}
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
		return
	}

	// first page
	for index := maxTableIndex; index >= 0 && len(list) < fetchLimit; index-- {
		if fixedTableIndex >= 0 && index != fixedTableIndex {
			continue
		}
		var tableData []model.ChainEvent
		q := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainEvent{BlockNum: uint(index) * model.SplitTableBlockNum}))
		q = q.Scopes(where...).Order("id desc").Limit(fetchLimit - len(list))
		if err := q.Find(&tableData).Error; err != nil {
			continue
		}
		list = append(list, tableData...)
	}
	hasNext = len(list) > limit
	if hasNext {
		list = list[:limit]
	}
	hasPrev = false
	return
}

func (d *Dao) GetEventsByIndex(extrinsicIndex string) []model.ChainEvent {
	var Event []model.ChainEvent
	indexArr := strings.Split(extrinsicIndex, "-")
	query := d.db.Scopes(model.TableNameFunc(model.ChainEvent{BlockNum: util.StringToUInt(indexArr[0])})).
		Where("extrinsic_index = ?", extrinsicIndex).Find(&Event)
	if query.Error != nil {
		return nil
	}
	return Event
}

func (d *Dao) GetEventByIdx(ctx context.Context, index string) *model.ChainEvent {
	var Event model.ChainEvent
	indexArr := strings.Split(index, "-")
	if len(indexArr) < 2 {
		return nil
	}
	query := d.db.WithContext(ctx).Scopes(model.TableNameFunc(model.ChainEvent{BlockNum: util.StringToUInt(indexArr[0])})).
		Where("block_num = ?", indexArr[0]).
		Where("event_idx = ?", indexArr[1]).First(&Event)
	if query == nil || query.Error != nil {
		return nil
	}
	return &Event
}
