package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"strings"
)

func (d *Dao) CreateEvent(txn *GormDB, event *model.ChainEvent) error {
	var incrCount int
	e := model.ChainEvent{
		ID:           event.Id(),
		EventIndex:   event.EventIndex,
		BlockNum:     event.BlockNum,
		ModuleId:     event.ModuleId,
		Params:       event.Params,
		EventIdx:     event.EventIdx,
		EventId:      event.EventId,
		ExtrinsicIdx: event.ExtrinsicIdx,
	}
	query := txn.Scopes(model.IgnoreDuplicate).Create(&e)
	if query.RowsAffected > 0 {
		incrCount++
		_ = d.IncrMetadata(context.TODO(), "count_event", incrCount)
	}
	return query.Error
}

func (d *Dao) GetEventByBlockNum(blockNum uint, where ...string) []model.ChainEventJson {
	var events []model.ChainEventJson
	queryOrigin := d.db.Model(model.ChainEvent{BlockNum: blockNum}).Where("block_num = ?", blockNum)
	for _, w := range where {
		queryOrigin = queryOrigin.Where(w)
	}
	query := queryOrigin.Order("id asc").Scan(&events)
	if query.Error != nil {
		return nil
	}
	return events
}

func (d *Dao) GetEventList(ctx context.Context, page, row int, order string, where ...string) ([]model.ChainEvent, int) {
	var Events []model.ChainEvent

	var count int64

	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	for index := blockNum / int(model.SplitTableBlockNum); index >= 0; index-- {
		var tableData []model.ChainEvent
		var tableCount int64
		fmt.Println("index", index)
		queryOrigin := d.db.WithContext(ctx).Model(model.ChainEvent{BlockNum: uint(index) * model.SplitTableBlockNum})
		for _, w := range where {
			queryOrigin = queryOrigin.Where(w)
		}

		queryOrigin.Count(&tableCount)

		if tableCount == 0 {
			continue
		}
		preCount := count
		count += tableCount
		if len(Events) >= row {
			continue
		}
		query := queryOrigin.Order(fmt.Sprintf("id %s", order)).Offset(page*row - int(preCount)).Limit(row - len(Events)).Scan(&tableData)
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
	query := d.db.Model(model.ChainEvent{BlockNum: util.StringToUInt(indexArr[0])}).
		Where("event_index = ?", extrinsicIndex).Scan(&Event)
	if query.Error != nil {
		return nil
	}
	return Event
}

func (d *Dao) GetEventByIdx(index string) *model.ChainEvent {
	var Event model.ChainEvent
	indexArr := strings.Split(index, "-")
	if len(indexArr) < 2 {
		return nil
	}
	query := d.db.Model(model.ChainEvent{BlockNum: util.StringToUInt(indexArr[0])}).
		Where("block_num = ?", indexArr[0]).
		Where("event_idx = ?", indexArr[1]).Scan(&Event)
	if query == nil || query.Error != nil {
		return nil
	}
	return &Event
}
