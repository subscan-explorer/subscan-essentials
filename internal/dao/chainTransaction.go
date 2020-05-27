package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
	"github.com/shopspring/decimal"
	"strings"
)

func (d *Dao) CreateTransaction(c context.Context, txn *GormDB, e *model.ChainExtrinsic, blockTimestamp int) error {
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
			amount = util.DecimalFromInterface(v.Value).Div(decimal.New(1, int32(substrate.BalanceAccuracy)))
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
		Fee:                e.Fee,
		Finalized:          e.Finalized,
	}
	query := txn.Create(&t)
	if query.RowsAffected == 0 {
		return errors.New("query rows not affected")
	}
	return query.Error
}

func (d *Dao) GetTransactionCount(c context.Context, where ...string) int {
	var count int
	blockNum, _ := d.GetFillAlreadyBlockNum(context.TODO())
	for index := blockNum / model.SplitTableBlockNum; index >= 0; index-- {
		var tableCount int
		queryOrigin := d.db.Model(model.ChainTransaction{BlockNum: index * model.SplitTableBlockNum})
		for _, w := range where {
			queryOrigin = queryOrigin.Where(w)
		}
		queryOrigin.Count(&tableCount)
		count += tableCount
	}
	return count
}

func (d *Dao) GetTransactionList(c context.Context, page, row int, order string, where ...string) ([]model.ExtrinsicsJson, int) {
	var txs []model.ExtrinsicsJson
	var count int

	blockNum, _ := d.GetFillAlreadyBlockNum(context.TODO())

	var transferQuery bool
	for index := blockNum / model.SplitTableBlockNum; index >= 0; index-- {
		var tableData []model.ExtrinsicsJson
		var tableCount int
		queryOrigin := d.db.Model(model.ChainTransaction{BlockNum: index * model.SplitTableBlockNum})

		for _, w := range where {
			queryOrigin = queryOrigin.Where(w)
			if strings.Contains(w, "transfer") {
				transferQuery = true
			}
		}

		if transferQuery {
			for _, w := range where {
				if strings.Contains(w, "from_hex") {
					queryOrigin = queryOrigin.Or(strings.Replace(w, "from_hex", "destination", 1))
					break
				}
			}
		}

		queryOrigin.Count(&tableCount)

		if tableCount == 0 {
			continue
		}
		preCount := count
		count += tableCount
		if len(txs) >= row {
			continue
		}
		query := queryOrigin.Order(fmt.Sprintf("block_num %s", order)).Offset(page*row - preCount).Limit(row - len(txs)).Scan(&tableData)
		if query == nil || query.Error != nil || query.RecordNotFound() {
			continue
		}
		txs = append(txs, tableData...)
	}
	return txs, count
}
