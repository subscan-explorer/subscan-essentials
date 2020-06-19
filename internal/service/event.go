package service

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/service/balances"
	"github.com/itering/subscan/internal/service/system"
	"github.com/itering/subscan/internal/util"
	"github.com/shopspring/decimal"
	"strings"
)

func (s *Service) GetEventList(page, row int, order string, where ...string) ([]model.ChainEvent, int) {
	c := context.TODO()
	return s.Dao.GetEventList(c, page, row, order, where...)
}

func (s *Service) GetEventByIndex(index string) []model.ChainEvent {
	return s.Dao.GetEventsByIndex(index)
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

	s.Dao.DropEventNotFinalizedData(blockNum, finalized)

	for _, event := range e {

		event.ModuleId = strings.ToLower(event.ModuleId)
		event.ExtrinsicHash = hashMap[fmt.Sprintf("%d-%d", blockNum, event.ExtrinsicIdx)]
		event.EventIndex = fmt.Sprintf("%d-%d", blockNum, event.ExtrinsicIdx)
		event.Finalized = finalized
		event.BlockNum = blockNum

		if err = s.Dao.CreateEvent(c, txn, &event); err == nil && finalized {
			go s.AnalysisEvent(blockHash, blockTimestamp, event, spec, feeMap[event.EventIndex])
		} else {
			return 0, err
		}
		if !util.StringInSlice(event.EventId, []string{"ExtrinsicSuccess", "ExtrinsicFailed"}) {
			eventCount++
		}

	}
	return eventCount, err
}

func (s *Service) AnalysisEvent(blockHash string, blockTimestamp int, event model.ChainEvent, spec int, fee decimal.Decimal) {
	paramEvent, err := model.ParsingEventParam(event.Params)
	if err != nil {
		return
	}
	switch event.ModuleId {
	case "system":
		system.EmitEvent(system.NewEvent(s.Dao, &event, paramEvent, blockHash, blockTimestamp, spec), event.EventId)

	case "balances", "kton": // ring
		balances.EmitEvent(balances.New(s.Dao, &event, paramEvent, blockTimestamp, fee), event.EventId)
	}

}
