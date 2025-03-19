package service

import (
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"strings"
)

func (s *Service) AddEvent(
	txn *dao.GormDB,
	block *model.ChainBlock,
	e []model.ChainEvent) (eventCount int, err error) {

	for index, event := range e {
		event.ModuleId = strings.ToLower(event.ModuleId)
		event.EventIndex = fmt.Sprintf("%d-%d", block.BlockNum, index)
		event.BlockNum = block.BlockNum
		if err = s.dao.CreateEvent(txn, &event); err == nil {
			go s.emitEvent(block, &event)
		} else {
			return 0, err
		}
		eventCount++
	}
	return eventCount, err
}
