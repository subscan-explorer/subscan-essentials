package service

import (
	"context"
	"encoding/json"
	"fmt"
	"subscan-end/internal/dao"
	"subscan-end/internal/model"
	"time"
)

func (s *Service) GetEventList(page, row int) (*[]model.ChainEventJson, int) {
	c := context.TODO()
	return s.dao.GetEventList(c, page, row)
}

func (s *Service) GetEventByIndex(index string) *[]model.ChainEventJson {
	c := context.TODO()
	return s.dao.GetEventsByIndex(c, index)
}

func (s *Service) AddEvent(c context.Context, txn *dao.GormDB, blockNum, blockTimestamp int, blockHash string, e []model.ChainEvent, hashMap map[string]string) int {
	var eventCount int
	for _, event := range e {
		if event.ModuleId != "system" {
			eventIndex := fmt.Sprintf("%d-%d", blockNum, event.ExtrinsicIdx)
			event.ExtrinsicHash = hashMap[eventIndex]
			_ = s.dao.CreateEvent(c, txn, blockNum, &event)
			eventCount += 1
			go s.AnalysisEvent(c, blockHash, blockNum, blockTimestamp, event)
		}
	}
	return eventCount
}

func (s *Service) AnalysisEvent(c context.Context, blockHash string, blockNum, blockTimestamp int, event model.ChainEvent) {
	var paramStruct []model.EventParam
	bj, _ := json.Marshal(event.Params.([]interface{}))
	if err := json.Unmarshal(bj, &paramStruct); err != nil {
		return
	}
	switch event.ModuleId {
	case "session":
		if event.EventId == "NewSession" {
			go s.SessionDeal(blockHash, uint(paramStruct[0].Value.(float64)), blockNum)
		}
	case "balances": //ring
		if event.EventId == "Transfer" {
			s.dao.UpdateAccountBalance(c, paramStruct[0].ValueRaw, event.ModuleId)
			s.dao.UpdateAccountBalance(c, paramStruct[1].ValueRaw, event.ModuleId)
			s.dao.IncrStatTransfer(c, time.Unix(int64(blockTimestamp), 0))
		}
	case "kton": //kton
		if event.EventId == "TokenTransfer" {
			s.dao.UpdateAccountBalance(c, paramStruct[0].ValueRaw, event.ModuleId)
			s.dao.UpdateAccountBalance(c, paramStruct[1].ValueRaw, event.ModuleId)
			s.dao.IncrStatTransfer(c, time.Unix(int64(blockTimestamp), 0))
		}
	case "indices":
		if event.EventId == "NewAccountIndex" {
			s.dao.UpdateAccountIndex(c, paramStruct[0].ValueRaw, int(paramStruct[1].Value.(float64)))
		}
	}
}
