package dao

import (
	"github.com/itering/subscan/model"
)

func (d *Dao) CreateLog(txn *GormDB, ce *model.ChainLog) error {
	query := txn.Scopes(d.TableNameFunc(ce), model.IgnoreDuplicate).Create(ce)
	return query.Error
}

func (d *Dao) GetLogByBlockNum(blockNum uint) []model.ChainLogJson {
	var logs []model.ChainLogJson
	query := d.db.Model(&model.ChainLog{BlockNum: blockNum}).
		Where("block_num =?", blockNum).Order("id asc").Scan(&logs)
	if query == nil || query.Error != nil {
		return nil
	}
	return logs
}
