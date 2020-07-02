package service

import (
	"context"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/lib/substrate"
	"github.com/itering/subscan/lib/substrate/storage"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"strings"
)

func (s *Service) GetLogByIndex(index string) *model.ChainLogJson {
	c := context.TODO()
	return s.dao.GetLogsByIndex(c, index)
}

func (s *Service) EmitLog(c context.Context, txn *dao.GormDB, blockNum int, l []storage.DecoderLog, validatorList []string, finalized bool) (validator string, err error) {
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
