package dao

import (
	"context"
	"fmt"
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

func (d *Dao) GetEventList(ctx context.Context, page, row int, order string, fixedTableIndex int, afterId uint, where ...model.Option) ([]model.ChainEvent, int) {
	var Events []model.ChainEvent

	var count int64

	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	maxTableIndex := blockNum / int(model.SplitTableBlockNum)
	if afterId > 0 {
		maxTableIndex = int(afterId/model.SplitTableBlockNum) / model.IdGenerateCoefficient
	}
	if fixedTableIndex >= 0 {
		maxTableIndex = fixedTableIndex
	}
	for index := maxTableIndex; index >= 0; index-- {
		var tableData []model.ChainEvent
		var tableCount int64

		if fixedTableIndex >= 0 && index != fixedTableIndex {
			continue
		}

		queryOrigin := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainEvent{BlockNum: uint(index) * model.SplitTableBlockNum}))
		queryOrigin.Scopes(where...)
		queryOrigin.Count(&tableCount)

		if tableCount == 0 {
			continue
		}
		preCount := count
		count += tableCount
		if len(Events) >= row {
			continue
		}
		if afterId > 0 {
			queryOrigin = queryOrigin.Where("id < ?", afterId)
		}
		query := queryOrigin.Order(fmt.Sprintf("id %s", order)).Offset(page*row - int(preCount)).Limit(row - len(Events)).Find(&tableData)
		if query == nil || query.Error != nil {
			continue
		}
		Events = append(Events, tableData...)

	}
	return Events, int(count)
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
