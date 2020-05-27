package service

import (
	"context"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/libs/substrate"
	"github.com/itering/subscan/libs/substrate/storage"
	"github.com/itering/subscan/util"
	"strings"
)

func (s *Service) GetLogList(page, row int) (*[]model.ChainLogJson, int) {
	c := context.TODO()
	return s.dao.GetLogList(c, page, row)
}

func (s *Service) GetLogByIndex(index string) *model.ChainLogJson {
	c := context.TODO()
	return s.dao.GetLogsByIndex(c, index)
}

func (s *Service) EmitLog(c context.Context, txn *dao.GormDB, blockNum int, l []storage.DecoderLog, validatorList []string, finalized bool) (validator string, err error) {
	s.dao.DropLogsNotFinalizedData(blockNum, finalized)
	for index, logData := range l {
		dataStr := util.InterfaceToString(logData.Value)

		if strings.ToLower(logData.Type) == "other" {
			dataStr = dealDarwiniaMmrRoot(dataStr)
		}

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

func dealDarwiniaMmrRoot(value string) string {
	if util.IsDarwinia {
		mmr := storage.MerkleMountainRangeRootLog{MmrRoot: strings.ReplaceAll(value, `"`, "")[8:]}
		return util.InterfaceToString(mmr)
	}
	return value
}
