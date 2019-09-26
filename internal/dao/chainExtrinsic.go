package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"subscan-end/internal/model"
	"subscan-end/utiles"
	"subscan-end/utiles/ss58"
)

func (d *Dao) CreateExtrinsic(c context.Context, txn *GormDB, blockNum string, index, blockTimestamp int, success bool, extrinsic *model.ChainExtrinsic) error {
	params, _ := json.Marshal(extrinsic.Params)
	extrinsicHash := ""
	if extrinsic.ExtrinsicHash != "" {
		extrinsicHash = utiles.AddHex(extrinsic.ExtrinsicHash)
	}
	ce := &model.ChainExtrinsic{
		BlockTimestamp:     blockTimestamp,
		ExtrinsicIndex:     fmt.Sprintf("%s-%d", utiles.HexToNumStr(blockNum), index),
		BlockNum:           utiles.StringToInt(utiles.HexToNumStr(blockNum)),
		ValueRaw:           extrinsic.ValueRaw,
		ExtrinsicLength:    extrinsic.ExtrinsicLength,
		VersionInfo:        extrinsic.VersionInfo,
		CallCode:           extrinsic.CallCode,
		CallModuleFunction: extrinsic.CallModuleFunction,
		CallModule:         extrinsic.CallModule,
		Params:             string(params),
		AccountLength:      extrinsic.AccountLength,
		AccountId:          extrinsic.AccountId,
		AccountIndex:       extrinsic.AccountIndex,
		Signature:          extrinsic.Signature,
		Era:                extrinsic.Era,
		ExtrinsicHash:      extrinsicHash,
		Nonce:              extrinsic.Nonce,
		Success:            success,
		IsSigned:           extrinsic.Signature != "",
	}
	query := txn.Create(&ce)
	if query.RowsAffected > 0 {
		_ = d.IncrMetadata(c, "count_extrinsic", 1)
	}
	if err := d.CreateTransaction(c, ce, blockTimestamp); err == nil {
		_ = d.IncrMetadata(c, "count_signed_extrinsic", 1)
	}
	return query.Error
}

func (d *Dao) GetExtrinsicsByBlockNum(c context.Context, blockNum int) *[]model.ChainExtrinsicJson {
	var extrinsics []model.ChainExtrinsicJson
	query := d.db.Model(model.ChainExtrinsic{}).Where("block_num = ?", blockNum).Order("id asc").Scan(&extrinsics)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &extrinsics
}

func (d *Dao) GetExtrinsicList(c context.Context, page, row int, order string, queryWhere ...string) (*[]model.ChainExtrinsicJson, int) {
	var extrinsics []model.ChainExtrinsicJson
	queryOrigin := d.db.Model(&model.ChainExtrinsic{})
	for _, w := range queryWhere {
		queryOrigin = queryOrigin.Where(w)
	}
	query := queryOrigin.Offset(page * row).Limit(row).Order(fmt.Sprintf("block_num %s", order)).Scan(&extrinsics)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return &extrinsics, 0
	}
	var count int
	if len(queryWhere) == 0 {
		m, _ := d.GetMetadata(c)
		count = utiles.StringToInt(m["count_extrinsic"])
	} else {
		queryOrigin.Count(&count)
	}
	return &extrinsics, count
}

func (d *Dao) GetExtrinsicsByHash(c context.Context, hash string) *model.ChainExtrinsicJson {
	var extrinsic model.ChainExtrinsicJson
	query := d.db.Model(model.ChainExtrinsic{}).Where("extrinsic_hash = ?", hash).Scan(&extrinsic)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &extrinsic
}

func (d *Dao) GetExtrinsicsByIndex(c context.Context, index string) *model.ChainExtrinsicJson {
	var extrinsic model.ChainExtrinsicJson
	query := d.db.Model(model.ChainExtrinsic{}).Where("extrinsic_index = ?", index).Scan(&extrinsic)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return &extrinsic
}

func (d *Dao) GetTimestamp(c context.Context, params []model.ExtrinsicParam) (timestamp int) {
	for _, p := range params {
		if p.Name == "now" {
			return int(p.Value.(float64))
		}
	}
	return
}

func (d *Dao) GetExtrinsicsDetailByHash(c context.Context, hash string) *model.ExtrinsicDetail {
	var extrinsic model.ChainExtrinsic
	query := d.db.Model(model.ChainExtrinsic{}).Where("extrinsic_hash = ?", hash).First(&extrinsic)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return d.extrinsicsAsDetail(c, &extrinsic)
}

func (d *Dao) GetExtrinsicsDetailByIndex(c context.Context, index string) *model.ExtrinsicDetail {
	var extrinsic model.ChainExtrinsic
	query := d.db.Model(model.ChainExtrinsic{}).Where("extrinsic_index = ?", index).Scan(&extrinsic)
	if query == nil || query.RecordNotFound() {
		return nil
	}
	return d.extrinsicsAsDetail(c, &extrinsic)
}

func (d *Dao) extrinsicsAsDetail(c context.Context, e *model.ChainExtrinsic) *model.ExtrinsicDetail {
	detail := model.ExtrinsicDetail{
		BlockTimestamp:     e.BlockTimestamp,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		BlockNum:           e.BlockNum,
		CallModule:         e.CallModule,
		CallModuleFunction: e.CallModuleFunction,
		AccountId:          ss58.Encode(e.AccountId),
		Signature:          e.Signature,
		Nonce:              e.Nonce,
		ExtrinsicHash:      e.ExtrinsicHash,
		Success:            e.Success,
	}
	var params []model.ExtrinsicParam
	_ = json.Unmarshal([]byte(e.Params.([]uint8)), &params)
	detail.Params = &params
	if detail.ExtrinsicHash != "" {
		detail.Event = d.GetEventsByIndex(c, e.ExtrinsicIndex)
	}
	if detail.CallModuleFunction == TransferModule {
		var dest string
		var amount decimal.Decimal
		for _, v := range params {
			if v.Type == "Address" {
				dest = v.Value.(string)
			}
			if v.Type == "Compact<Balance>" {
				amount = utiles.FloatToDecimal(v.Value.(float64))
			}
		}
		t := model.TransferJson{
			From:    detail.AccountId,
			To:      ss58.Encode(dest),
			Module:  detail.CallModule,
			Hash:    detail.ExtrinsicHash,
			Amount:  amount,
			Success: detail.Success,
		}
		detail.Transfer = &t
	}
	return &detail
}
