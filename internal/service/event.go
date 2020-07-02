package service

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	"github.com/shopspring/decimal"
	"strings"
)

func (s *Service) GetEventList(page, row int, order string, where ...string) ([]model.ChainEvent, int) {
	c := context.TODO()
	return s.dao.GetEventList(c, page, row, order, where...)
}

func (s *Service) GetEventByIndex(index string) []model.ChainEvent {
	return s.dao.GetEventsByIndex(index)
}

func (s *Service) AddEvent(
	c context.Context,
	txn *dao.GormDB,
	blockNum, blockTimestamp int,
	blockHash string,
	e []model.ChainEvent,
	hashMap map[string]string,
	finalized bool,
	spec int,
	feeMap map[string]decimal.Decimal,
) (eventCount int, err error) {
	s.dao.DropEventNotFinalizedData(blockNum, finalized)

	for _, event := range e {
		event.ModuleId = strings.ToLower(event.ModuleId)
		event.ExtrinsicHash = hashMap[fmt.Sprintf("%d-%d", blockNum, event.ExtrinsicIdx)]
		event.EventIndex = fmt.Sprintf("%d-%d", blockNum, event.ExtrinsicIdx)
		event.Finalized = finalized
		event.BlockNum = blockNum

		if err = s.dao.CreateEvent(c, txn, &event); err == nil && finalized {
			go s.afterEvent(spec, blockTimestamp, blockHash, &event, feeMap[event.EventIndex])
		} else {
			return 0, err
		}
		eventCount++
	}
	return eventCount, err
}

func (s *Service) afterEvent(spec, blockTimestamp int, blockHash string, event *model.ChainEvent, fee decimal.Decimal) {
	for _, plugin := range plugins.RegisteredPlugins {
		p := plugin()
		_ = p.ProcessEvent(spec, blockTimestamp, blockHash, event, fee)
	}

}
