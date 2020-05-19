package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"subscan-end/internal/dao"
	"subscan-end/internal/model"
	"subscan-end/libs/substrate"
	"subscan-end/libs/substrate/protos/codec_protos"
	"subscan-end/utiles"
)

func (s *Service) GetLogList(page, row int) (*[]model.ChainLogJson, int) {
	c := context.TODO()
	return s.dao.GetLogList(c, page, row)
}

func (s *Service) GetLogByIndex(index string) *model.ChainLogJson {
	c := context.TODO()
	return s.dao.GetLogsByIndex(c, index)
}

func (s *Service) EmitLog(c context.Context, txn *dao.GormDB, blockHash, blockNum, decodeLog string) (validator string) {
	var l []codec_protos.DecoderLog
	_ = json.Unmarshal([]byte(decodeLog), &l)
	for index, logData := range l {
		data, _ := json.Marshal(logData.Value)
		s.dao.CreateLog(c, txn, blockNum, index, logData, data)
		if logData.Index == "PreRuntime" {
			c, _, err := websocket.DefaultDialer.Dial(utiles.ProviderEndPoint, nil)
			if err != nil {
				return ""
			}
			validatorList, _ := s.GetValidatorFromSub(c, blockHash) //todo,websocket connection not init
			validator = substrate.ExtractAuthor(data, validatorList)
			c.Close()
		}
	}
	return validator
}
