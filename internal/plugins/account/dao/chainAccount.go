package dao

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/internal/plugins/account/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/rpc"
	"github.com/itering/subscan/internal/util"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func TouchAccount(db *gorm.DB, address string) (*model.ChainAccount, error) {
	var account model.ChainAccount
	address = util.TrimHex(address)
	query := db.FirstOrCreate(&account, &model.ChainAccount{Address: address})
	return &account, query.Error
}

func updateBalance(db *gorm.DB, account *model.ChainAccount, balance decimal.Decimal) error {
	u := map[string]interface{}{"balance": balance}
	query := db.Model(account).Update(u)
	if query == nil || query.Error != nil || query.RowsAffected == 0 {
		return errors.New("update balance account error")
	}
	return nil
}

func GetBalanceFromNetwork(db *gorm.DB, address string) (decimal.Decimal, error) {
	balance, _, err := rpc.GetFreeBalance(nil, address, "")
	if err != nil {
		log.Error("GetBalanceFromNetwork error %v", err)
		return decimal.Zero, err
	}
	return balance.Div(decimal.New(1, int32(substrate.BalanceAccuracy))), nil
}

func UpdateAccountBalance(db *gorm.DB, account *model.ChainAccount) (decimal.Decimal, error) {
	balance, err := GetBalanceFromNetwork(db, account.Address)
	if err == nil {
		_ = updateBalance(db, account, balance)
	}
	return balance, err
}

func ResetAccountNonce(db *gorm.DB, address string) {
	account, err := TouchAccount(db, address)
	if err != nil {
		return
	}
	_ = db.Model(account).Update(model.ChainAccount{Nonce: 0})
}

func UpdateAccountLock(db *gorm.DB, address, currency string) error {
	balance, err := rpc.GetAccountLock(nil, address, currency)
	if err != nil {
		log.Error("UpdateAccountLock err %v", err)
		return err
	}
	u := map[string]interface{}{"lock": balance}
	query := db.Model(model.ChainAccount{}).Where("address = ?", address).Update(u)
	if query == nil || query.Error != nil || query.RowsAffected == 0 {
		return errors.New("update balance lock error")
	}
	return nil
}

func GetAccountList(db *gorm.DB, page, row int, order, field string, queryWhere ...string) ([]*model.ChainAccount, int) {
	var accounts []*model.ChainAccount
	queryOrigin := db.Model(&model.ChainAccount{})
	if field == "" {
		field = "id"
	}
	for _, w := range queryWhere {
		queryOrigin = queryOrigin.Where(w)
	}
	query := queryOrigin.Order(fmt.Sprintf("%s %s", field, order)).Offset(page * row).Limit(row).Scan(&accounts)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return accounts, 0
	}
	var count int
	queryOrigin.Count(&count)
	return accounts, count
}
