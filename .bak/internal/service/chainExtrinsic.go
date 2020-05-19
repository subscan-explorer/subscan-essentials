package service

import (
	"context"
	"encoding/json"
	"fmt"
	"subscan-end/internal/dao"
	"subscan-end/internal/model"
	"subscan-end/libs/substrate"
	"subscan-end/utiles"
)

func (s *Service) GetExtrinsicList(page, row int, order string, query ...string) (*[]model.ChainExtrinsicJson, int) {
	c := context.TODO()
	return s.dao.GetExtrinsicList(c, page, row, order, query...)
}

func (s *Service) GetExtrinsicDetailByHash(hash string) *model.ExtrinsicDetail {
	c := context.TODO()
	return s.dao.GetExtrinsicsDetailByHash(c, hash)
}

func (s *Service) GetExtrinsicByHash(hash string) *model.ChainExtrinsicJson {
	c := context.TODO()
	return s.dao.GetExtrinsicsByHash(c, hash)
}

func (s *Service) GetExtrinsicByIndex(index string) *model.ExtrinsicDetail {
	c := context.TODO()
	return s.dao.GetExtrinsicsDetailByIndex(c, index)
}

func (s *Service) GetTransferList(page, row int) *[]model.TransferJson {
	c := context.TODO()
	return s.dao.GetTransfers(c, page, row)
}

func (s *Service) createExtrinsic(c context.Context, txn *dao.GormDB, blockNum, decodeExtrinsics string, successMap map[string]bool) (int, int, map[string]string) {
	var (
		blockTimestamp int
		e              []model.ChainExtrinsic
	)
	_ = json.Unmarshal([]byte(decodeExtrinsics), &e)
	hash := make(map[string]string)
	for index, extrinsic := range e {
		params, _ := json.Marshal(extrinsic.Params)
		var paramsInstant []model.ExtrinsicParam
		_ = json.Unmarshal(params, &paramsInstant)
		if extrinsic.CallModule == "timestamp" {
			blockTimestamp = s.dao.GetTimestamp(c, paramsInstant)
		}
		var successFlag, ok bool
		if successFlag, ok = successMap[fmt.Sprintf("%s-%d", utiles.HexToNumStr(blockNum), index)]; ok == false {
			successFlag = true
		}
		extrinsicIndex := fmt.Sprintf("%s-%d", utiles.HexToNumStr(blockNum), index)
		hash[extrinsicIndex] = extrinsic.ExtrinsicHash
		if err := s.dao.CreateExtrinsic(c, txn, blockNum, index, blockTimestamp, successFlag, &extrinsic); err == nil && successFlag && extrinsic.ExtrinsicHash != "" {
			go s.AnalysisExtrinsic(c, &extrinsic, paramsInstant)
		}
	}
	return len(e), blockTimestamp, hash
}

func (s *Service) AnalysisExtrinsic(c context.Context, e *model.ChainExtrinsic, params []model.ExtrinsicParam) {
	switch e.CallModule {
	case "staking":
		s.extrinsicStaking(c, e, params)
	case "session":
		s.extrinsicSession(c, e, params)
	}
}

func (s *Service) extrinsicSession(c context.Context, e *model.ChainExtrinsic, params []model.ExtrinsicParam) {
	if len(params) == 0 {
		return
	}
	switch e.CallModuleFunction {
	case "set_keys": // Keys, proof
		keys := params[0].Value.(map[string]interface{})
		_ = s.dao.CreateValidatorSessionOrUpdate(c, e.AccountId, keys["col1"].(string))
	}
}

func (s *Service) extrinsicStaking(c context.Context, e *model.ChainExtrinsic, params []model.ExtrinsicParam) {
	if len(params) == 0 {
		return
	}
	// need rpc to find out
	switch e.CallModuleFunction {
	case "validate": // name, ratio, unstake_threshold
		if len(params) == 3 && utiles.NetworkNode == utiles.DarwiniaNetwork {
			nodeName := params[0].Value.(string)
			_ = s.dao.UpdateValidatorName(c, nodeName, e.AccountId)
		}
	case "set_controller": // Keys, proof
		controller := params[0].Value.(string)
		_ = s.dao.CreateValidatorStashOrUpdate(c, controller, e.AccountId)
	case "bond": //controller, value, payee, promise_month
		controller := params[0].Value.(string)
		_ = s.dao.CreateValidatorStashOrUpdate(c, controller, e.AccountId)
	case "set_payee": //payee
		rewardDestination := substrate.GetRewardDestinationEnum()
		index := utiles.StringToInt(params[0].ValueRaw)
		if index >= len(rewardDestination) {
			return
		}
		_ = s.dao.SetRewardAccount(c, e.AccountId, rewardDestination[index])
	}
}
