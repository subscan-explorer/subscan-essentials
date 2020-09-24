package dao

import (
	"fmt"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/shopspring/decimal"
)

func GetAccountList(db storage.DB, page, row int) ([]model.Account, int) {
	var accounts []model.Account
	db.FindBy(&accounts, nil, &storage.Option{PluginPrefix: "balance"})
	return accounts, len(accounts)
}

func NewAccount(db storage.DB, accountId string) error {
	accountId = util.TrimHex(accountId)
	err := db.Create(&model.Account{Address: accountId})
	if err != nil {
		err = AfterAccountCreate(db, accountId)
	}
	return err
}

func AfterAccountCreate(db storage.DB, accountId string) error {
	accountDataRaw, err := rpc.ReadStorage(nil, "system", "account", "", accountId)
	if err != nil {
		return err
	}
	accountData := new(model.AccountData)
	accountDataRaw.ToAny(accountData)
	return db.Update(model.Account{}, fmt.Sprintf("address = '%s'", accountId), map[string]interface{}{
		"nonce":   accountData.Nonce,
		"balance": accountData.Data.Free.Add(accountData.Data.Reserved),
		"lock":    decimal.Max(accountData.Data.MiscFrozen, accountData.Data.FeeFrozen),
	})
}
