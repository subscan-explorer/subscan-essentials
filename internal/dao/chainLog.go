package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/libs/substrate/storage"
	"github.com/itering/subscan/util"
	"strings"
)

func (d *Dao) CreateLog(c context.Context, txn *GormDB, blockNum, index int, logData *storage.DecoderLog, data []byte, finalized bool) error {
	ce := model.ChainLog{
		LogIndex:   fmt.Sprintf("%d-%d", blockNum, index),
		BlockNum:   blockNum,
		LogType:    logData.Type,
		OriginType: logData.Type,
		Data:       string(data),
		Finalized:  finalized,
	}
	query := txn.Create(&ce)
	return d.checkDBError(query.Error)
}

func (d *Dao) DropLogsNotFinalizedData(blockNum int, finalized bool) bool {
	var delExist bool
	if finalized {
		query := d.db.Where("block_num = ?", blockNum).
			Delete(model.ChainLog{BlockNum: blockNum})
		delExist = query.RowsAffected > 0
	}
	return delExist
}

// TODO
func (d *Dao) GetLogList(c context.Context, page, row int) (*[]model.ChainLogJson, int) {
	var Logs []model.ChainLogJson
	query := d.db.Model(&model.ChainLog{}).Offset(page * row).Limit(row).Order("block_num desc").Scan(&Logs)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return &Logs, 0
	}
	var count int
	d.db.Model(&model.ChainLog{}).Count(&count)
	return &Logs, count
}

func (d *Dao) GetLogsByIndex(c context.Context, index string) *model.ChainLogJson {
	var Log model.ChainLogJson
	indexArr := strings.Split(index, "-")
	query := d.db.Model(model.ChainLog{BlockNum: util.StringToInt(indexArr[0])}).Where("log_index = ?", index).Scan(&Log)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &Log
}

func (d *Dao) GetLogByBlockNum(c context.Context, blockNum int) *[]model.ChainLogJson {
	var logs []model.ChainLogJson
	query := d.db.Model(&model.ChainLog{BlockNum: blockNum}).
		Where("block_num =?", blockNum).Order("id asc").Scan(&logs)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return &logs
}
