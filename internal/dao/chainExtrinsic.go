package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc"
	"strings"
)

func (d *Dao) CreateExtrinsic(c context.Context, txn *GormDB, extrinsic []model.ChainExtrinsic, signedExtrinsicCount int) error {
	if len(extrinsic) == 0 {
		return nil
	}
	query := txn.Scopes(d.TableNameFunc(&extrinsic[0]), model.IgnoreDuplicate).CreateInBatches(extrinsic, 2000)
	if query.RowsAffected > 0 {
		_ = d.IncrMetadata(c, "count_extrinsic", int(query.RowsAffected))
		_ = d.IncrMetadata(c, "count_signed_extrinsic", signedExtrinsicCount)
	}
	return query.Error
}

func (d *Dao) GetExtrinsicCount(ctx context.Context, queryWhere ...model.Option) int64 {
	var count int64
	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	for index := blockNum / int(model.SplitTableBlockNum); index >= 0; index-- {
		var tableDataCount int64
		q := d.db.WithContext(ctx).Scopes(d.TableNameFunc(&model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum}))
		q = q.Scopes(queryWhere...)
		q.Model(&model.ChainExtrinsic{}).Count(&tableDataCount)
		count += tableDataCount
	}
	return count

}

func (d *Dao) GetAccountExtrinsicMapping(ctx context.Context, accountId string) []int {
	var mapping model.AccountExtrinsicMapping
	query := d.db.WithContext(ctx).Where("account_id = ?", accountId).First(&mapping)
	if query != nil && query.Error == nil {
		return mapping.ExtrinsicTable
	}
	return nil
}

// GetExtrinsicListCursor implements bidirectional cursor pagination using id as cursor.
// When afterId > 0, fetch records with id < afterId in DESC order.
// When beforeId > 0, fetch records with id > beforeId in ASC order then reverse.
func (d *Dao) GetExtrinsicListCursor(c context.Context, limit int, fixedTableIndex int, beforeId, afterId uint, accountId string, queryWhere ...model.Option) (list []model.ChainExtrinsic, hasPrev, hasNext bool) {
	fetchLimit := limit + 1
	blockNum, _ := d.GetFillBestBlockNum(context.TODO())
	maxTableIndex := blockNum / int(model.SplitTableBlockNum)
	if afterId > 0 {
		maxTableIndex = int(afterId/model.SplitTableBlockNum) / model.IdGenerateCoefficient
	}
	if fixedTableIndex >= 0 {
		maxTableIndex = fixedTableIndex
	}

	var accountExtrinsics []int
	if accountId != "" {
		// find extrinsic table by AccountExtrinsicMapping
		accountExtrinsics = d.GetAccountExtrinsicMapping(c, accountId)
	}

	var checkTableIndex = func(index int) bool {
		if len(accountId) == 0 || (len(accountId) > 0 && util.IntInSlice(index, accountExtrinsics)) {
			return true
		}
		return false
	}

	if afterId > 0 { // next page
		for index := maxTableIndex; index >= 0 && len(list) < fetchLimit; index-- {
			if (fixedTableIndex >= 0 && index != fixedTableIndex) || !checkTableIndex(index) {
				continue
			}
			var tableData []model.ChainExtrinsic
			q := d.db.WithContext(c).Scopes(d.TableNameFunc(&model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum}))
			q = q.Scopes(queryWhere...)
			q = q.Where("id < ?", afterId).Order("id desc").Limit(fetchLimit - len(list))
			if err := q.Find(&tableData).Error; err != nil {
				continue
			}
			list = append(list, tableData...)
		}
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = true
		return
	}

	if beforeId > 0 { // previous page
		startIdx := int(beforeId/model.SplitTableBlockNum) / model.IdGenerateCoefficient
		if fixedTableIndex >= 0 {
			startIdx = fixedTableIndex
		}
		for index := startIdx; index <= maxTableIndex && len(list) < fetchLimit; index++ {
			if (fixedTableIndex >= 0 && index != fixedTableIndex) || !checkTableIndex(index) {
				continue
			}
			var tableData []model.ChainExtrinsic
			q := d.db.WithContext(c).Scopes(d.TableNameFunc(&model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum}))
			q = q.Scopes(queryWhere...)
			if index == startIdx {
				q = q.Where("id > ?", beforeId)
			}
			q = q.Order("id asc").Limit(fetchLimit - len(list))
			if err := q.Find(&tableData).Error; err != nil {
				continue
			}
			list = append(list, tableData...)
		}
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		// reverse to keep DESC order in response
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
		return
	}

	// first page
	for index := maxTableIndex; index >= 0 && len(list) < fetchLimit; index-- {
		if (fixedTableIndex >= 0 && index != fixedTableIndex) || !checkTableIndex(index) {
			continue
		}
		var tableData []model.ChainExtrinsic
		q := d.db.WithContext(c).Scopes(d.TableNameFunc(&model.ChainExtrinsic{BlockNum: uint(index) * model.SplitTableBlockNum}))
		q = q.Scopes(queryWhere...).Order("id desc").Limit(fetchLimit - len(list))
		if err := q.Find(&tableData).Error; err != nil {
			continue
		}
		list = append(list, tableData...)
	}
	hasNext = len(list) > limit
	if hasNext {
		list = list[:limit]
	}
	hasPrev = false
	return
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
