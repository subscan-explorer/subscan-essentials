package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/plugins"
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

			blockTimestamp = s.getTimestamp(extrinsic.Params)
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

	s.dao.DropExtrinsicNotFinalizedData(c, blockNum, finalized)

	for _, extrinsic := range extrinsicList {
		if err = s.dao.CreateExtrinsic(c, txn, &extrinsic); err == nil {
			go s.afterExtrinsic(spec, &extrinsic, eventMap[extrinsic.ExtrinsicIndex])
		} else {
			return 0, 0, nil, nil, err
		}
	}
	return len(e), blockTimestamp, hash, extrinsicFee, err
}

func (s *Service) getTimestamp(param interface{}) (timestamp int) {
	var paramsInstant []model.ExtrinsicParam
	util.UnmarshalToAnything(&paramsInstant, param)
	for _, p := range paramsInstant {
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

func (s *Service) afterExtrinsic(spec int, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) {
	for _, plugin := range plugins.RegisteredPlugins {
		p := plugin()
		_ = p.ProcessExtrinsic(spec, extrinsic, events)
	}
}
