package service

import (
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"strings"
)

func (s *Service) AddEvent(txn *dao.GormDB, block *model.ChainBlock, e []model.ChainEvent) (err error) {
	for _, event := range e {
		event.ModuleId = strings.ToLower(event.ModuleId)
		event.ExtrinsicIndex = fmt.Sprintf("%d-%d", block.BlockNum, event.ExtrinsicIdx)
		event.BlockNum = block.BlockNum
		if err = s.dao.CreateEvent(txn, &event); err == nil {
			go s.emitEvent(block, &event)
		} else {
			return err
		}
	}
	return err
}
