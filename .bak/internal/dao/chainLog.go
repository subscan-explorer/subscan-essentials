package dao

import (
	"context"
	"fmt"
	"subscan-end/internal/model"
	"subscan-end/libs/substrate/protos/codec_protos"
	"subscan-end/utiles"
)

func (d *Dao) CreateLog(c context.Context, txn *GormDB, blockNum string, index int, logData codec_protos.DecoderLog, data []byte) {
	ce := &model.ChainLog{
		LogIndex:   fmt.Sprintf("%s-%d", utiles.HexToNumStr(blockNum), index),
		BlockNum:   utiles.StringToInt(utiles.HexToNumStr(blockNum)),
		LogType:    logData.Index,
		OriginType: logData.Type,
		Data:       string(data),
	}
	txn.Create(&ce)
	return
}

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
	query := d.db.Model(model.ChainLog{}).Where("log_index = ?", index).Scan(&Log)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &Log
}

func (d *Dao) GetLogByBlockNum(c context.Context, blockNum int) *[]model.ChainLogJson {
	var logs []model.ChainLogJson
	query := d.db.Model(&model.ChainLog{}).Where("block_num = ?", blockNum).Order("id asc").Scan(&logs)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil
	}
	return &logs
}
