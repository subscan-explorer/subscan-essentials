package dao

import (
	"context"
	"strings"
	"subscan-end/internal/model"
	"subscan-end/utiles/ss58"
)

func (d *Dao) GetTransfers(c context.Context, page, row int, where ...string) *[]model.TransferJson {
	var txs []model.ChainTransaction
	query := d.db.Model(model.ChainTransaction{})
	for _, w := range where {
		query = query.Where(w)
		if strings.Contains(w, "from_hex") {
			query = query.Or(strings.Replace(w, "from_hex", "destination", 1))
		}
	}
	query = query.Order("block_num desc").Offset(page*row).Limit(row).Where("`call_module_function` = ?", TransferModule).Scan(&txs)
	tj := []model.TransferJson{}
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return &tj
	}
	for _, t := range txs {
		tj = append(tj, model.TransferJson{
			From:           ss58.Encode(t.FromHex),
			To:             ss58.Encode(t.Destination),
			Module:         t.CallModule,
			Hash:           t.Hash,
			Amount:         t.Amount,
			BlockTimestamp: t.BlockTimestamp,
			ExtrinsicIndex: t.ExtrinsicIndex,
			BlockNum:       t.BlockNum,
			Success:        t.Success,
		})
	}

	return &tj
}

func (d *Dao) GetTransferCount(c context.Context, where ...string) int {
	var count int
	query := d.db.Model(model.ChainTransaction{})
	for _, w := range where {
		query = query.Where(w)
		if strings.Contains(w, "from_hex") {
			query = query.Or(strings.Replace(w, "from_hex", "destination", 1))
		}
	}
	query.Where("`call_module_function` = ?", TransferModule).Count(&count)
	return count
}
