package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"subscan-end/internal/model"
	"subscan-end/utiles"
)

func (d *Dao) CreateEvent(c context.Context, txn *GormDB, blockNum int, event *model.ChainEvent) error {
	var incrCount int
	params, _ := json.Marshal(event.Params)
	query := txn.Create(&model.ChainEvent{
		EventIndex:    fmt.Sprintf("%d-%d", blockNum, event.ExtrinsicIdx),
		BlockNum:      blockNum,
		Type:          event.Type,
		ModuleId:      event.ModuleId,
		Params:        string(params),
		Phase:         event.Phase,
		EventIdx:      event.EventIdx,
		EventId:       event.EventId,
		ExtrinsicIdx:  event.ExtrinsicIdx,
		ExtrinsicHash: event.ExtrinsicHash,
	})
	if query.RowsAffected > 0 {
		incrCount += 1
	}
	_ = d.IncrMetadata(c, "count_event", incrCount)
	return query.Error
}

func (d *Dao) GetEventByBlockNum(c context.Context, blockNum int) *[]model.ChainEventJson {
	var events []model.ChainEventJson
	query := d.db.Model(model.ChainEvent{}).Where("block_num = ?", blockNum).Order("id asc").Scan(&events)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &events
}

func (d *Dao) GetEventList(c context.Context, page, row int) (*[]model.ChainEventJson, int) {
	var Events []model.ChainEventJson
	query := d.db.Model(&model.ChainEvent{}).Offset(page * row).Limit(row).Order("block_num desc").Order("id asc").Scan(&Events)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return &Events, 0
	}
	m, _ := d.GetMetadata(c)
	return &Events, utiles.StringToInt(m["count_event"])
}

func (d *Dao) GetEventsByIndex(c context.Context, index string) *[]model.ChainEventJson {
	var Event []model.ChainEventJson
	query := d.db.Model(model.ChainEvent{}).Where("event_index = ?", index).Scan(&Event)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &Event
}
