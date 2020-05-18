package dao

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"subscan-end/internal/model"
	"subscan-end/utiles"
	"subscan-end/utiles/ss58"
)

func (d *Dao) CreateTransaction(c context.Context, e *model.ChainExtrinsic, blockTimestamp int) error {
	var (
		params []model.ExtrinsicParam
		dest   string
		amount decimal.Decimal
	)
	_ = json.Unmarshal([]byte(e.Params.(string)), &params)

	if e.ExtrinsicHash == "" {
		return errors.New("no ExtrinsicHash")
	}
	for _, v := range params {
		if v.Type == "Address" {
			dest = v.Value.(string)
		}
		if v.Type == "Compact<Balance>" {
			amount = utiles.FloatToDecimal(v.Value.(float64))
		}
	}
	t := model.ChainTransaction{
		Signature:          e.Signature,
		ExtrinsicIndex:     e.ExtrinsicIndex,
		FromHex:            e.AccountId,
		Destination:        dest,
		Hash:               e.ExtrinsicHash,
		BlockNum:           e.BlockNum,
		BlockTimestamp:     e.BlockTimestamp,
		CallCode:           e.CallCode,
		CallModuleFunction: e.CallModuleFunction,
		CallModule:         e.CallModule,
		Params:             e.Params,
		Success:            e.Success,
		Amount:             amount,
	}
	query := d.db.Create(&t)
	if query.RowsAffected == 0 {
		return errors.New("query rows not affected")
	}
	d.UpdateAccountCountExtrinsic(c, e.AccountId)
	return query.Error
}

func (d *Dao) GetTransactionByAccount(c context.Context, address string, page, row int) (*[]model.ChainTransactionJson, int) {
	var txs []model.ChainTransactionJson
	query := d.db.Model(model.ChainTransaction{}).Offset(page*row).Limit(row).Where("`from_hex` = ?", address).Scan(&txs)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil, 0
	}
	for i, tx := range txs {
		txs[i].Destination = ss58.Encode(tx.Destination)
		txs[i].FromHex = ss58.Encode(tx.FromHex)
	}
	var count int
	d.db.Model(model.ChainTransaction{}).Where("`from_hex` = ?", address).Count(&count)
	return &txs, count
}
