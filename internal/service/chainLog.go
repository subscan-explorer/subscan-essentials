package service

import (
	"context"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"strings"
)

func (s *Service) EmitLog(c context.Context, txn *dao.GormDB, blockHash string, blockNum int, l []storage.DecoderLog, finalized bool) (validator string, err error) {
	validatorList, _ := rpc.GetValidatorFromSub(nil, blockHash)
	s.dao.DropLogsNotFinalizedData(blockNum, finalized)
	for index, logData := range l {
		dataStr := util.InterfaceToString(logData.Value)

		if err = s.dao.CreateLog(c, txn, blockNum, index, &logData, []byte(dataStr), finalized); err != nil {
			return "", err
		}

		// check validator
		if strings.EqualFold(logData.Type, "PreRuntime") {
			validator = substrate.ExtractAuthor([]byte(dataStr), validatorList)
		}

	}
	return validator, err
}
