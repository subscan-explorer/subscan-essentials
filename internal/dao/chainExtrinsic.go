package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc"
	"strings"
)

func (d *Dao) CreateExtrinsic(c context.Context, txn *GormDB, extrinsic *model.ChainExtrinsic) error {
	ce := model.ChainExtrinsic{
		ID:                 extrinsic.ID,
		BlockTimestamp:     extrinsic.BlockTimestamp,
		ExtrinsicIndex:     extrinsic.ExtrinsicIndex,
		BlockNum:           extrinsic.BlockNum,
		CallModuleFunction: extrinsic.CallModuleFunction,
		CallModule:         extrinsic.CallModule,
		Params:             extrinsic.Params,
		AccountId:          extrinsic.AccountId,
		Signature:          extrinsic.Signature,
		Era:                extrinsic.Era,
		ExtrinsicHash:      util.AddHex(extrinsic.ExtrinsicHash),
		Nonce:              extrinsic.Nonce,
		Success:            extrinsic.Success,
		IsSigned:           extrinsic.Signature != "",
		Fee:                extrinsic.Fee,
	}
	query := txn.Scopes(d.TableNameFunc(&ce), model.IgnoreDuplicate).Scopes(model.IgnoreDuplicate).Create(&ce)
	if query.RowsAffected > 0 {
		_ = d.IncrMetadata(c, "count_extrinsic", 1)
		if ce.IsSigned {
			_ = d.IncrMetadata(c, "count_signed_extrinsic", 1)
		}
	}
	return query.Error
}

func (d *Dao) GetExtrinsicsByBlockNum(blockNum uint) []model.ChainExtrinsicJson {
	var extrinsics []model.ChainExtrinsic
	query := d.db.Model(model.ChainExtrinsic{BlockNum: blockNum}).
		Where("block_num = ?", blockNum).Order("id asc").Scan(&extrinsics)
	if query == nil || query.Error != nil {
		return nil
	}
	var list []model.ChainExtrinsicJson
	for _, extrinsic := range extrinsics {
		list = append(list, *d.ExtrinsicsAsJson(&extrinsic))
	}
	return list
}

func (d *Dao) GetExtrinsicList(c context.Context, page, row int, _ string, queryWhere ...model.Option) ([]model.ChainExtrinsic, int) {
	var extrinsics []model.ChainExtrinsic
	var count int64

	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	for index := blockNum / int(model.SplitTableBlockNum); index >= 0; index-- {
		var (
			tableData  []model.ChainExtrinsic
			tableCount int64
		)

		queryOrigin := d.db.Scopes(d.TableNameFunc(&model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum}))

		queryOrigin.Scopes(queryWhere...)
		queryOrigin.Count(&tableCount)

		if tableCount == 0 {
			continue
		}
		preCount := count
		count += tableCount
		if len(extrinsics) >= row {
			continue
		}
		query := queryOrigin.Order("id desc").Offset(page*row - int(preCount)).Limit(row - len(extrinsics)).Scan(&tableData)
		if query == nil || query.Error != nil {
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

func (d *Dao) GetExtrinsicsByHash(c context.Context, hash string) *model.ChainExtrinsic {
	var extrinsic model.ChainExtrinsic
	blockNum, _ := d.GetFillBestBlockNum(c)
	for index := blockNum / int(model.SplitTableBlockNum); index >= 0; index-- {
		query := d.db.Model(model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum}).Where("extrinsic_hash = ?", hash).Order("id asc").Limit(1).Scan(&extrinsic)
		if query != nil && query.Error == nil {
			return &extrinsic
		}
	}
	return nil
}

func (d *Dao) GetExtrinsicsDetailByHash(c context.Context, hash string) *model.ExtrinsicDetail {
	if extrinsic := d.GetExtrinsicsByHash(c, hash); extrinsic != nil {
		return d.extrinsicsAsDetail(c, extrinsic)
	}
	return nil
}

func (d *Dao) GetExtrinsicsDetailByIndex(c context.Context, index string) *model.ExtrinsicDetail {
	var extrinsic model.ChainExtrinsic
	indexArr := strings.Split(index, "-")
	query := d.db.Model(model.ChainExtrinsic{BlockNum: util.StringToUInt(indexArr[0])}).
		Where("extrinsic_index = ?", index).Scan(&extrinsic)
	if query == nil || query.Error != nil {
		return nil
	}
	return d.extrinsicsAsDetail(c, &extrinsic)
}

func (d *Dao) extrinsicsAsDetail(ctx context.Context, e *model.ChainExtrinsic) *model.ExtrinsicDetail {
	detail := model.ExtrinsicDetail{
		BlockTimestamp:     e.BlockTimestamp,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		BlockNum:           e.BlockNum,
		CallModule:         e.CallModule,
		CallModuleFunction: e.CallModuleFunction,
		AccountId:          address.Encode(e.AccountId),
		Signature:          e.Signature,
		Nonce:              e.Nonce,
		ExtrinsicHash:      e.ExtrinsicHash,
		Success:            e.Success,
		Fee:                e.Fee,
		Finalized:          true,
	}
	d.FindLifeTime(ctx, &detail, e.Era)
	return &detail
}

func (d *Dao) ExtrinsicsAsJson(e *model.ChainExtrinsic) *model.ChainExtrinsicJson {
	ej := &model.ChainExtrinsicJson{
		Id:                 e.ID,
		BlockNum:           e.BlockNum,
		BlockTimestamp:     e.BlockTimestamp,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		ExtrinsicHash:      e.ExtrinsicHash,
		Success:            e.Success,
		CallModule:         e.CallModule,
		CallModuleFunction: e.CallModuleFunction,
		Params:             e.Params,
		AccountId:          address.Encode(e.AccountId),
		Signature:          e.Signature,
		Nonce:              e.Nonce,
		Fee:                e.Fee,
	}
	return ej
}

func (d *Dao) FindLifeTime(_ context.Context, detail *model.ExtrinsicDetail, era string) {
	if detail.Signature == "" {
		return
	}
	if mortal := substrate.DecodeMortal(era); mortal != nil {
		detail.Lifetime = &model.Lifetime{
			Birth: mortal.Birth(uint64(detail.BlockNum)),
			Death: mortal.Death(uint64(detail.BlockNum)),
		}
	}
}
