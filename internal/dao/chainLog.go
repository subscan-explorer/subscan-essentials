package dao

import (
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"strings"
)

func (d *Dao) CreateLog(txn *GormDB, ce *model.ChainLog) error {
	query := txn.Create(ce)
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

func (d *Dao) GetLogsByIndex(index string) *model.ChainLogJson {
	var Log model.ChainLogJson
	indexArr := strings.Split(index, "-")
	query := d.db.Model(model.ChainLog{BlockNum: util.StringToInt(indexArr[0])}).Where("log_index = ?", index).Scan(&Log)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &Log
}

func (d *Dao) GetLogByBlockNum(blockNum int) []model.ChainLogJson {
	var logs []model.ChainLogJson
	query := d.db.Model(&model.ChainLog{BlockNum: blockNum}).
		Where("block_num =?", blockNum).Order("id asc").Scan(&logs)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return logs
}
