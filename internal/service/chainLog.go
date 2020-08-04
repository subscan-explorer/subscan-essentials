package service

import (
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"strings"
)

func (s *Service) EmitLog(txn *dao.GormDB, blockHash string, blockNum int, l []storage.DecoderLog, finalized bool, validatorList []string) (validator string, err error) {
	s.dao.DropLogsNotFinalizedData(blockNum, finalized)
	for index, logData := range l {
		dataStr := util.ToString(logData.Value)

		ce := model.ChainLog{
			LogIndex:  fmt.Sprintf("%d-%d", blockNum, index),
			BlockNum:  blockNum,
			LogType:   logData.Type,
			Data:      dataStr,
			Finalized: finalized,
		}
		if err = s.dao.CreateLog(txn, &ce); err != nil {
			return "", err
		}

		// check validator
		if strings.EqualFold(ce.LogType, "PreRuntime") {
			validator = substrate.ExtractAuthor([]byte(dataStr), validatorList)
		}

	}
	return validator, err
}
