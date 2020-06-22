package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/service/transaction"
	"github.com/itering/subscan/internal/util"
	"github.com/shopspring/decimal"
	"strings"
)

func (s *Service) createExtrinsic(c context.Context,
	txn *dao.GormDB,
	blockNum int,
	encodeExtrinsics []string,
	decodeExtrinsics []map[string]interface{},
	eventMap map[string][]model.ChainEvent,
	finalized bool,
	spec int,
) (int, int, map[string]string, map[string]decimal.Decimal, error) {

	var (
		blockTimestamp int
		e              []model.ChainExtrinsic
		extrinsicList  []model.ChainExtrinsic
		err            error
	)
	extrinsicFee := make(map[string]decimal.Decimal)
	eb, _ := json.Marshal(decodeExtrinsics)
	_ = json.Unmarshal(eb, &e)
	hash := make(map[string]string)

	for index, extrinsic := range e {
		extrinsic.CallModule = strings.ToLower(extrinsic.CallModule)

		if extrinsic.CallModule == "timestamp" {
			var paramsInstant = model.ParsingExtrinsicParam(extrinsic.Params)
			blockTimestamp = s.getTimestamp(c, paramsInstant)
		}

		extrinsic.BlockNum = blockNum
		extrinsic.ExtrinsicIndex = fmt.Sprintf("%d-%d", extrinsic.BlockNum, index)

		extrinsic.Success = s.getExtrinsicSuccess(eventMap[extrinsic.ExtrinsicIndex])

		extrinsic.Finalized = finalized

		hash[extrinsic.ExtrinsicIndex] = extrinsic.ExtrinsicHash
		extrinsic.BlockTimestamp = blockTimestamp

		if extrinsic.ExtrinsicHash != "" {
			extrinsic.Fee = transaction.GetExtrinsicFee(encodeExtrinsics[index])
			extrinsicFee[extrinsic.ExtrinsicIndex] = extrinsic.Fee
		}

		extrinsicList = append(extrinsicList, extrinsic)
	}

	s.Dao.DropExtrinsicNotFinalizedData(c, blockNum, finalized)

	for _, extrinsic := range extrinsicList {
		nonce := s.getAccountNewNonce(extrinsic.AccountId, extrinsic.Nonce, eventMap[extrinsic.ExtrinsicIndex]) // check ReapedAccount

		extrinsicValue := extrinsic
		err = s.Dao.CreateExtrinsic(c, txn, &extrinsicValue, nonce)
		if err != nil {
			return 0, 0, nil, nil, err
		}
	}
	return len(e), blockTimestamp, hash, extrinsicFee, err
}

func (s *Service) getTimestamp(c context.Context, params []model.ExtrinsicParam) (timestamp int) {
	for _, p := range params {
		if p.Name == "now" {
			return util.IntFromInterface(p.Value)
		}
	}
	return
}

func (s *Service) getExtrinsicSuccess(e []model.ChainEvent) bool {
	for _, event := range e {
		if strings.ToLower(event.ModuleId) == "system" {
			return strings.ToLower(event.EventId) != "extrinsicfailed"
		}
	}
	return true
}

func (s *Service) getAccountNewNonce(accountID string, extrinsicNonce int, e []model.ChainEvent) int {
	for _, event := range e {
		if event.EventId == "ReapedAccount" {
			if params, err := model.ParsingEventParam(event.Params); err == nil && len(params) == 2 {
				if util.TrimHex(util.InterfaceToString(params[0].Value)) == accountID {
					return 0
				}
			}
		}
	}
	return extrinsicNonce + 1
}
