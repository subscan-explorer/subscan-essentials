package dao

import (
	"context"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/model"
	bModel "github.com/itering/subscan/plugins/balance/model"
	"gorm.io/gorm"
)

func CreateTransfer(ctx context.Context, d *Storage, transfer *bModel.Transfer) error {
	db := d.Dao.GetDbInstance().(*gorm.DB)
	query := db.WithContext(ctx).Scopes(model.IgnoreDuplicate).Create(transfer)
	if query.RowsAffected > 0 {
		_, _ = d.Pool.HINCRBY(ctx, model.MetadataCacheKey(), "total_transfer", 1)
		_ = RefreshAccount(ctx, d, model.CheckoutParamValueAddress(transfer.Sender))
		_ = RefreshAccount(ctx, d, model.CheckoutParamValueAddress(transfer.Receiver))
	}
	return query.Error
}

func TransfersCursor(ctx context.Context, db storage.DB, limit int, before, after *uint, opts ...model.Option) ([]bModel.Transfer, bool, bool) {
	var list []bModel.Transfer
	d := db.GetDbInstance().(*gorm.DB)
	fetch := limit + 1
	q := d.WithContext(ctx).Model(bModel.Transfer{}).Scopes(opts...)
	var hasPrev, hasNext bool
	if after != nil && *after > 0 {
		q = q.Where("id < ?", *after).Order("id desc")
	} else if before != nil && *before > 0 {
		q = q.Where("id > ?", *before).Order("id asc")
	} else {
		q = q.Order("id desc")
	}
	q = q.Limit(fetch).Find(&list)
	if q.Error != nil {
		return nil, false, false
	}
	if before != nil && *before > 0 {
		hasPrev = len(list) > limit
		if hasPrev {
			list = list[:limit]
		}
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
		hasNext = true
	} else {
		hasNext = len(list) > limit
		if hasNext {
			list = list[:limit]
		}
		hasPrev = after != nil && *after > 0
	}
	return list, hasPrev, hasNext
}
