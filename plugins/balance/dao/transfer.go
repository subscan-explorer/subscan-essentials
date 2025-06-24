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

func Transfers(ctx context.Context, db storage.DB, page model.Option, opts ...model.Option) ([]bModel.Transfer, int) {
	var list []bModel.Transfer
	d := db.GetDbInstance().(*gorm.DB)
	var count int64
	q := d.WithContext(ctx).Model(bModel.Transfer{}).Scopes(opts...).Count(&count)
	if q.Error != nil {
		return nil, 0
	}
	q = d.WithContext(ctx).Model(bModel.Transfer{}).Scopes(page).Scopes(opts...).Order("id desc").Find(&list)
	return list, int(count)
}
