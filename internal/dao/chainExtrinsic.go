package dao

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
)

func (d *Dao) CreateExtrinsic(c context.Context, txn *GormDB, extrinsic *model.ChainExtrinsic) error {
	existing, err := findOne[model.ChainExtrinsic](&d.ReadOnlyDao, "id", where("extrinsic_index = ?", extrinsic.ExtrinsicIndex), nil)
	if err != nil {
		return err
	}
	if existing != nil {
		extrinsic.ID = existing.ID
	}
	ce := model.ChainExtrinsic{
		ID:                 extrinsic.ID,
		BlockTimestamp:     extrinsic.BlockTimestamp,
		ExtrinsicIndex:     extrinsic.ExtrinsicIndex,
		BlockNum:           extrinsic.BlockNum,
		ExtrinsicLength:    extrinsic.ExtrinsicLength,
		VersionInfo:        extrinsic.VersionInfo,
		CallCode:           extrinsic.CallCode,
		CallModuleFunction: extrinsic.CallModuleFunction,
		CallModule:         extrinsic.CallModule,
		Params:             util.ToString(extrinsic.Params),
		AccountId:          extrinsic.AccountId,
		Signature:          extrinsic.Signature,
		Era:                extrinsic.Era,
		ExtrinsicHash:      util.AddHex(extrinsic.ExtrinsicHash),
		Nonce:              extrinsic.Nonce,
		Success:            extrinsic.Success,
		IsSigned:           extrinsic.Signature != "",
		Fee:                extrinsic.Fee,
	}
	query := txn.Save(&ce)
	if query.RowsAffected > 0 {
		_ = d.IncrMetadata(c, "count_extrinsic", 1)
		if ce.IsSigned {
			_ = d.IncrMetadata(c, "count_signed_extrinsic", 1)
		}
	}
	return d.checkDBError(query.Error)
}

func (d *ReadOnlyDao) GetExtrinsicsByBlockNum(blockNum int) []model.ChainExtrinsicJson {
	var extrinsics []model.ChainExtrinsic
	query := d.db.Model(model.ChainExtrinsic{BlockNum: blockNum}).
		Where("block_num = ?", blockNum).Order("id asc").Scan(&extrinsics)
	if query == nil || RecordNotFound(query) {
		return nil
	}
	var list []model.ChainExtrinsicJson
	for _, extrinsic := range extrinsics {
		list = append(list, *d.ExtrinsicsAsJson(&extrinsic))
	}
	return list
}

func (d *ReadOnlyDao) GetExtrinsicList(c context.Context, page, row int, order string, queryWhere ...string) ([]model.ChainExtrinsic, int) {
	var extrinsics []model.ChainExtrinsic
	var count int64

	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	for index := blockNum / model.SplitTableBlockNum; index >= 0; index-- {
		var tableData []model.ChainExtrinsic
		var tableCount int64
		queryOrigin := d.db.Model(model.ChainExtrinsic{BlockNum: index * model.SplitTableBlockNum})
		for _, w := range queryWhere {
			queryOrigin = queryOrigin.Where(w)
		}

		queryOrigin.Count(&tableCount)

		if tableCount == 0 {
			continue
		}
		preCount := count
		count += tableCount
		if len(extrinsics) >= row {
			continue
		}
		query := queryOrigin.Order("block_num desc").Offset(page*row - int(preCount)).Limit(row - len(extrinsics)).Scan(&tableData)
		if query == nil || query.Error != nil || RecordNotFound(query) {
			continue
		}
		extrinsics = append(extrinsics, tableData...)

	}

	if len(queryWhere) == 0 {
		m, _ := d.GetMetadata(c)
		count = int64(util.StringToInt(m["count_extrinsic"]))
	}
	return extrinsics, int(count)
}

func (d *ReadOnlyDao) GetExtrinsicsByHash(c context.Context, hash string) *model.ChainExtrinsic {
	var extrinsic model.ChainExtrinsic
	blockNum, _ := d.GetFillBestBlockNum(c)
	for index := blockNum / (model.SplitTableBlockNum); index >= 0; index-- {
		query := d.db.Model(model.ChainExtrinsic{BlockNum: index * model.SplitTableBlockNum}).Select("*").Where("extrinsic_hash = ?", hash).Order("id asc").Limit(1).Scan(&extrinsic)
		if query != nil && !RecordNotFound(query) {
			if extrinsic.Params != nil {
				extrinsic.Params = *(extrinsic.Params.(*interface{}))
			}
			return &extrinsic
		}
	}
	return nil
}

func (d *ReadOnlyDao) GetExtrinsicsDetailByHash(c context.Context, hash string) *model.ExtrinsicDetail {
	if extrinsic := d.GetExtrinsicsByHash(c, hash); extrinsic != nil {
		return d.extrinsicsAsDetail(c, extrinsic)
	}
	return nil
}

func (d *ReadOnlyDao) GetExtrinsicsDetailByIndex(c context.Context, index string) *model.ExtrinsicDetail {
	var extrinsics []model.ChainExtrinsic
	indexArr := strings.Split(index, "-")
	query := d.db.Model(model.ChainExtrinsic{BlockNum: util.StringToInt(indexArr[0])}).
		Select("*").
		Where("extrinsic_index = ?", index).
		Limit(1).
		Find(&extrinsics)
	if query == nil || query.Error != nil || len(extrinsics) == 0 {
		return nil
	}
	extrinsic := extrinsics[0]
	if extrinsic.Params != nil {
		extrinsic.Params = *(extrinsic.Params.(*interface{}))
	}
	return d.extrinsicsAsDetail(c, &extrinsic)
}

func (d *ReadOnlyDao) extrinsicsAsDetail(c context.Context, e *model.ChainExtrinsic) *model.ExtrinsicDetail {
	detail := model.ExtrinsicDetail{
		BlockTimestamp:     e.BlockTimestamp,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		BlockNum:           e.BlockNum,
		CallModule:         e.CallModule,
		CallModuleFunction: e.CallModuleFunction,
		AccountId:          address.SS58AddressFromHex(e.AccountId),
		Signature:          e.Signature,
		Nonce:              e.Nonce,
		ExtrinsicHash:      e.ExtrinsicHash,
		Success:            e.Success,
		Fee:                e.Fee,
	}
	util.UnmarshalAny(&detail.Params, e.Params)

	if block := d.GetBlockByNum(detail.BlockNum); block != nil {
		detail.Finalized = block.Finalized
	}

	events := d.GetEventsByIndex(e.ExtrinsicIndex)
	for k, event := range events {
		events[k].Params = util.ToString(event.Params)
	}

	detail.Event = &events

	return &detail
}

func (d *ReadOnlyDao) ExtrinsicsAsJson(e *model.ChainExtrinsic) *model.ChainExtrinsicJson {
	ej := &model.ChainExtrinsicJson{
		BlockNum:           e.BlockNum,
		BlockTimestamp:     e.BlockTimestamp,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		ExtrinsicHash:      e.ExtrinsicHash,
		Success:            e.Success,
		CallModule:         e.CallModule,
		CallModuleFunction: e.CallModuleFunction,
		Params:             util.ToString(e.Params),
		AccountId:          address.SS58AddressFromHex(e.AccountId),
		Signature:          e.Signature,
		Nonce:              e.Nonce,
		Fee:                e.Fee,
	}
	var paramsInstant []model.ExtrinsicParam
	if err := json.Unmarshal([]byte(ej.Params), &paramsInstant); err != nil {
		for pi, param := range paramsInstant {
			if paramsInstant[pi].Type == "Address" {
				paramsInstant[pi].Value = address.SS58AddressFromHex(param.Value.(string))
			}
		}
		bp, _ := json.Marshal(paramsInstant)
		ej.Params = string(bp)
	}
	return ej
}
