package dao

import (
	"context"
	"github.com/itering/subscan/model"
)

func (d *Dao) CreateLog(txn *GormDB, ce []model.ChainLog) error {
	if len(ce) == 0 {
		return nil
	}
	query := txn.Scopes(d.TableNameFunc(ce[0]), model.IgnoreDuplicate).CreateInBatches(ce, 2000)
	return query.Error
}

func (d *Dao) GetLogByBlockNum(ctx context.Context, blockNum uint) []model.ChainLogJson {
	var logs []model.ChainLogJson
	query := d.db.WithContext(ctx).Scopes(model.TableNameFunc(&model.ChainLog{BlockNum: blockNum})).
		Where("block_num =?", blockNum).Order("id asc").Scan(&logs)
	if query == nil || query.Error != nil {
		return nil
	}
	return logs
}
