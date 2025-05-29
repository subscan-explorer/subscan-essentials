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
		UsedFee:            extrinsic.UsedFee,
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

func (d *Dao) GetExtrinsicList(c context.Context, page, row int, _ string, queryWhere ...model.Option) ([]model.ChainExtrinsic, int) {
	var extrinsics []model.ChainExtrinsic
	var count int64

	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	for index := blockNum / int(model.SplitTableBlockNum); index >= 0; index-- {
		var (
			tableData  []model.ChainExtrinsic
			tableCount int64
		)

		queryOrigin := d.db.WithContext(c).Scopes(d.TableNameFunc(&model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum}))

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
	return extrinsics, int(count)
}

func (d *Dao) GetExtrinsicsByHash(c context.Context, hash string) *model.ChainExtrinsic {
	var extrinsic model.ChainExtrinsic
	blockNum, _ := d.GetFillBestBlockNum(c)
	for index := blockNum / int(model.SplitTableBlockNum); index >= 0; index-- {
		query := d.db.WithContext(c).Scopes(model.TableNameFunc(model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum})).Where("extrinsic_hash = ?", hash).First(&extrinsic)
		if query != nil && query.Error == nil {
			return &extrinsic
		}
	}
	return nil
}

func (d *Dao) GetExtrinsicsByIndex(c context.Context, index string) *model.ChainExtrinsic {
	var extrinsic model.ChainExtrinsic
	indexArr := strings.Split(index, "-")
	query := d.db.WithContext(c).Scopes(model.TableNameFunc(model.ChainExtrinsic{BlockNum: util.StringToUInt(indexArr[0])})).
		Where("extrinsic_index = ?", index).Find(&extrinsic)
	if query == nil || query.Error != nil {
		return nil
	}
	return &extrinsic
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
	query := d.db.WithContext(c).Scopes(model.TableNameFunc(model.ChainExtrinsic{BlockNum: util.StringToUInt(indexArr[0])})).
		Where("extrinsic_index = ?", index).Find(&extrinsic)
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
		Params:             e.Params,
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
