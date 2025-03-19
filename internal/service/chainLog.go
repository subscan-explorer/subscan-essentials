package service

import (
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"strings"
)

func (s *Service) EmitLog(txn *dao.GormDB, blockNum uint, l []storage.DecoderLog, finalized bool, validatorList []string) (validator string, err error) {
	for index, logData := range l {
		var jsonRaw model.LogData
		switch v := logData.Value.(type) {
		case map[string]interface{}:
			jsonRaw = v
		default:
			jsonRaw = map[string]interface{}{"data": v}
		}

		ce := model.ChainLog{
			LogIndex:  fmt.Sprintf("%d-%d", blockNum, index),
			BlockNum:  blockNum,
			LogType:   logData.Type,
			Data:      jsonRaw,
			Finalized: finalized,
		}
		ce.ID = ce.Id()
		if err = s.dao.CreateLog(txn, &ce); err != nil {
			return "", err
		}

		// check validator
		if strings.EqualFold(ce.LogType, "PreRuntime") {
			validator = substrate.ExtractAuthor(jsonRaw.Bytes(), validatorList)
		}

	}
	return validator, err
}
