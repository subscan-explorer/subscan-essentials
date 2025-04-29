package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transfer struct {
	Id             uint            `json:"id" gorm:"primary_key;autoIncrement:false"`
	Sender         string          `json:"sender" gorm:"size:255;index:query_function"`
	Receiver       string          `json:"receiver" gorm:"size:255;index:query_function"`
	Amount         decimal.Decimal `json:"amount" gorm:"decimal(65)"`
	BlockTimestamp int64           `json:"block_timestamp" `
	Symbol         string          `json:"symbol" gorm:"size:255"`
	TokenId        string          `json:"token_id" gorm:"size:255"`
	ExtrinsicIndex string          `json:"extrinsic_index" gorm:"size:255;index:extrinsic_index"`
}

func CreateTransfer(ctx context.Context, d *Storage, transfer *Transfer) error {
	db := d.Dao.GetDbInstance().(*gorm.DB)
	query := db.WithContext(ctx).Scopes(model.IgnoreDuplicate).Create(transfer)
	if query.RowsAffected > 0 {
		_, _ = d.Pool.HINCRBY(ctx, model.MetadataCacheKey(), "total_transfer", 1)
		_ = RefreshAccount(ctx, d, model.CheckoutParamValueAddress(transfer.Sender))
		_ = RefreshAccount(ctx, d, model.CheckoutParamValueAddress(transfer.Receiver))
	}
	return query.Error
}
