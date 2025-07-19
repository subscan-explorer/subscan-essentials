package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
)

func (s *Service) createExtrinsic(ctx context.Context,
	txn *dao.GormDB,
	block *model.ChainBlock,
	extrinsics []model.ChainExtrinsic,
	encodeExtrinsics []string,
	eventMap map[string][]model.ChainEvent,
) (err error) {

	for index, extrinsic := range extrinsics {
		extrinsic.BlockNum = block.BlockNum
		extrinsic.ExtrinsicIndex = fmt.Sprintf("%d-%d", extrinsic.BlockNum, index)
		extrinsic.ID = extrinsic.Id()
		extrinsic.Success = s.getExtrinsicSuccess(eventMap[extrinsic.ExtrinsicIndex])
		extrinsic.BlockTimestamp = block.BlockTimestamp
		extrinsic.AccountId = address.Format(extrinsic.AccountId)
		if extrinsic.Signature != "" {
			weight, actualFee, isV2Weight := model.CheckoutWeight(eventMap[extrinsic.ExtrinsicIndex])
			extrinsic.Fee, extrinsic.UsedFee, err = GetExtrinsicFee(ctx, encodeExtrinsics[index], block.ParentHash, block.SpecVersion, weight, actualFee, isV2Weight)
			if err != nil {
				util.Logger().Error(fmt.Errorf("extrinsic %s GetExtrinsicFee err %v", extrinsic.ExtrinsicIndex, err))
			}
		}

		if err = s.dao.CreateExtrinsic(ctx, txn, &extrinsic); err == nil {
			if err = s.emitExtrinsic(ctx, block, &extrinsic); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (s *Service) ExtrinsicsAsJson(e *model.ChainExtrinsic) *model.ChainExtrinsicJson {
	ej := &model.ChainExtrinsicJson{
		BlockNum:           e.BlockNum,
		BlockTimestamp:     e.BlockTimestamp,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		ExtrinsicHash:      e.ExtrinsicHash,
		Success:            e.Success,
		CallModule:         e.CallModule,
		CallModuleFunction: e.CallModuleFunction,
		Params:             e.Params,
		AccountId:          address.Encode(e.AccountId),
		Signature:          e.Signature,
		Nonce:              e.Nonce,
		Fee:                e.Fee,
	}
	return ej
}

func FindOutBlockTime(extrinsics []model.ChainExtrinsic) int {
	for _, extrinsic := range extrinsics {
		if strings.EqualFold(extrinsic.CallModule, "timestamp") {
			params := model.ParsingExtrinsicParam(extrinsic.Params)
			for _, p := range params {
				if strings.EqualFold(p.Name, "now") {
					if strings.EqualFold(p.Type, "compact<U64>") {
						return int(util.Int64FromInterface(p.Value) / 1000)
					}
					return util.IntFromInterface(p.Value)
				}
			}
		}
	}
	return 0
}

func (s *Service) getExtrinsicSuccess(e []model.ChainEvent) bool {
	for _, event := range e {
		if strings.EqualFold(event.ModuleId, "system") {
			if !strings.EqualFold(event.EventId, "ExtrinsicSuccess") || !strings.EqualFold(event.EventId, "ExtrinsicFailed") {
				continue
			} else {
				return strings.EqualFold(event.EventId, "ExtrinsicSuccess")
			}
		}
	}
	return true
}
