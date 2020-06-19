package service

import (
	"context"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/storage"
	"github.com/itering/subscan/internal/util"
	"strings"
)

func (s *Service) GetLogList(page, row int) (*[]model.ChainLogJson, int) {
	c := context.TODO()
	return s.Dao.GetLogList(c, page, row)
}

func (s *Service) GetLogByIndex(index string) *model.ChainLogJson {
	c := context.TODO()
	return s.Dao.GetLogsByIndex(c, index)
}

func (s *Service) EmitLog(c context.Context, txn *dao.GormDB, blockNum int, l []storage.DecoderLog, validatorList []string, finalized bool) (validator string, err error) {
	s.Dao.DropLogsNotFinalizedData(blockNum, finalized)
	for index, logData := range l {
		dataStr := util.InterfaceToString(logData.Value)

		if err = s.Dao.CreateLog(c, txn, blockNum, index, &logData, []byte(dataStr), finalized); err != nil {
			return "", err
		}

		// check validator
		if strings.EqualFold(logData.Type, "PreRuntime") {
			validator = substrate.ExtractAuthor([]byte(dataStr), validatorList)
		}

	}
	return validator, err
}
