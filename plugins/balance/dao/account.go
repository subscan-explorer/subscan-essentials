package dao

import (
	"context"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func GetAccountList(db storage.DB, page, row int) ([]model.Account, int) {
	var accounts []model.Account
	opt := storage.Option{PluginPrefix: "balance", Page: page, PageSize: row}
	db.FindBy(&accounts, nil, &opt)
	return accounts, len(accounts)
}

func RefreshAccount(ctx context.Context, s storage.DB, accountId string) error {
	accountId = address.Format(accountId)
	db := s.GetDbInstance().(*gorm.DB)
	var account model.Account
	db.WithContext(ctx).Where("address = ?", accountId).FirstOrCreate(&account)
	return AfterAccountCreate(ctx, db, &account)
}

func AfterAccountCreate(ctx context.Context, db *gorm.DB, account *model.Account) error {
	accountDataRaw, err := rpc.ReadStorage(nil, "system", "account", "", account.Address)
	if err != nil {
		return err
	}
	accountData := new(model.AccountData)
	accountDataRaw.ToAny(accountData)

	return db.WithContext(ctx).Where("address = ?", account.Address).UpdateColumns(map[string]interface{}{
		"nonce":   accountData.Nonce,
		"balance": accountData.Data.Free.Add(accountData.Data.Reserved),
		"lock":    decimal.Max(accountData.Data.MiscFrozen, accountData.Data.FeeFrozen),
	}).Error
}
