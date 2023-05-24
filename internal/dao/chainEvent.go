package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
)

func (d *Dao) CreateEvent(txn *GormDB, event *model.ChainEvent) error {
	var incrCount int
	extrinsicHash := util.AddHex(event.ExtrinsicHash)
	e := model.ChainEvent{
		EventIndex:    event.EventIndex,
		BlockNum:      event.BlockNum,
		Type:          event.Type,
		ModuleId:      event.ModuleId,
		Params:        util.ToString(event.Params),
		EventIdx:      event.EventIdx,
		EventId:       event.EventId,
		ExtrinsicIdx:  event.ExtrinsicIdx,
		ExtrinsicHash: extrinsicHash,
	}
	query := txn.Create(&e)
	if query.RowsAffected > 0 {
		incrCount++
		_ = d.IncrMetadata(context.TODO(), "count_event", incrCount)
	}
	return d.checkDBError(query.Error)
}

func (d *ReadOnlyDao) GetEventByBlockNum(blockNum int, where ...string) []model.ChainEventJson {
	var events []model.ChainEventJson
	queryOrigin := d.db.Model(model.ChainEvent{BlockNum: blockNum}).Where("block_num = ?", blockNum)
	for _, w := range where {
		queryOrigin = queryOrigin.Where(w)
	}
	query := queryOrigin.Order("id asc").Scan(&events)
	if query == nil || RecordNotFound(query) {
		return nil
	}
	return events
}

func (d *ReadOnlyDao) GetEventList(page, row int, order string, where ...string) ([]model.ChainEvent, int) {
	var events []model.ChainEvent
	q := d.db.Model(&model.ChainEvent{})
	for _, w := range where {
		q = q.Where(w)
	}
	q.Order(fmt.Sprintf("block_num %s", order)).Offset(page * row).Limit(row).Scan(&events)

	return events, len(events)
}

func (d *ReadOnlyDao) GetEventsByIndex(extrinsicIndex string) []model.ChainEvent {
	var Event []model.ChainEvent
	indexArr := strings.Split(extrinsicIndex, "-")
	query := d.db.Model(model.ChainEvent{BlockNum: util.StringToInt(indexArr[0])}).
		Where("event_index = ?", extrinsicIndex).Scan(&Event)
	if query == nil || RecordNotFound(query) {
		return nil
	}
	return Event
}

func (d *ReadOnlyDao) GetEventByIdx(index string) *model.ChainEvent {
	var Event model.ChainEvent
	indexArr := strings.Split(index, "-")
	if len(indexArr) < 2 {
		return nil
	}
	query := d.db.Model(model.ChainEvent{BlockNum: util.StringToInt(indexArr[0])}).
		Where("block_num = ?", indexArr[0]).
		Where("event_idx = ?", indexArr[1]).Scan(&Event)
	if query == nil || RecordNotFound(query) {
		return nil
	}
	return &Event
}
