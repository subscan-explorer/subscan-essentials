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

func (d *Dao) DropEventNotFinalizedData(blockNum int, finalized bool) bool {
	var delExist bool
	if finalized {
		query := d.db.Where("block_num = ?", blockNum).Delete(model.ChainEvent{BlockNum: blockNum})
		delExist = query.RowsAffected > 0
	}
	return delExist
}

func (d *Dao) GetEventByBlockNum(blockNum int, where ...string) []model.ChainEventJson {
	var events []model.ChainEventJson
	queryOrigin := d.db.Model(model.ChainEvent{BlockNum: blockNum}).Where("block_num = ?", blockNum)
	for _, w := range where {
		queryOrigin = queryOrigin.Where(w)
	}
	query := queryOrigin.Order("id asc").Scan(&events)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return events
}

func (d *Dao) GetEventList(page, row int, order string, where ...string) ([]model.ChainEvent, int) {
	var Events []model.ChainEvent

	var count int

	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	for index := blockNum / model.SplitTableBlockNum; index >= 0; index-- {
		var tableData []model.ChainEvent
		var tableCount int
		queryOrigin := d.db.Model(model.ChainEvent{BlockNum: index * model.SplitTableBlockNum})
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
		query := queryOrigin.Order(fmt.Sprintf("block_num %s", order)).Offset(page*row - preCount).Limit(row - len(Events)).Scan(&tableData)
		if query == nil || query.Error != nil || query.RecordNotFound() {
			continue
		}
		Events = append(Events, tableData...)

	}
	return Events, count
}

func (d *Dao) GetEventsByIndex(extrinsicIndex string) []model.ChainEvent {
	var Event []model.ChainEvent
	indexArr := strings.Split(extrinsicIndex, "-")
	query := d.db.Model(model.ChainEvent{BlockNum: util.StringToInt(indexArr[0])}).
		Where("event_index = ?", extrinsicIndex).Scan(&Event)
	if query == nil || query.RecordNotFound() {
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
	query := d.db.Model(model.ChainEvent{BlockNum: util.StringToInt(indexArr[0])}).
		Where("block_num = ?", indexArr[0]).
		Where("event_idx = ?", indexArr[1]).Scan(&Event)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &Event
}
