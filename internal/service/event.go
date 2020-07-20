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

func (s *Service) AddEvent(
	c context.Context,
	txn *dao.GormDB,
	block *model.ChainBlock,
	e []model.ChainEvent,
	hashMap map[string]string,
	finalized bool,
	spec int,
	feeMap map[string]decimal.Decimal,
) (eventCount int, err error) {

	s.dao.DropEventNotFinalizedData(block.BlockNum, finalized)
	for _, event := range e {
		event.ModuleId = strings.ToLower(event.ModuleId)
		event.ExtrinsicHash = hashMap[fmt.Sprintf("%d-%d", block.BlockNum, event.ExtrinsicIdx)]
		event.EventIndex = fmt.Sprintf("%d-%d", block.BlockNum, event.ExtrinsicIdx)
		event.Finalized = finalized
		event.BlockNum = block.BlockNum

		if err = s.dao.CreateEvent(c, txn, &event); err == nil {
			go s.afterEvent(*block, event, feeMap[event.EventIndex])
		} else {
			return 0, err
		}
		eventCount++
	}
	return eventCount, err
}

func (s *Service) afterEvent(block model.ChainBlock, event model.ChainEvent, fee decimal.Decimal) {
	pBlock := block.AsPluginBlock()
	pEvent := event.AsPluginEvent()
	for _, plugin := range plugins.RegisteredPlugins {
		_ = plugin.ProcessEvent(pBlock, pEvent, fee)
	}

}
