package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/rpc"
	"github.com/itering/subscan/internal/util"
	"github.com/itering/subscan/internal/util/ss58"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (d *Dao) TouchAccount(c context.Context, address string) (*model.ChainAccount, error) {
	var account model.ChainAccount
	var err error
	address = util.TrimHex(address)
	query := d.db.First(&account, &model.ChainAccount{Address: address})
	if query == nil || query.Error != nil || query.RecordNotFound() {
		account = model.ChainAccount{Address: address, AccountIndex: -1}
		if query = d.db.Create(&account); query != nil {
			err = query.Error
		}
	}
	return &account, err
}

func (d *Dao) FindByAddress(address string) (*model.ChainAccount, error) {
	address = util.TrimHex(address)
	var account model.ChainAccount
	query := d.db.First(&account, &model.ChainAccount{Address: address})
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil, errors.New("touch account error")
	}
	return &account, nil
}

func (d *Dao) updateBalance(c context.Context, account *model.ChainAccount, balance, additional decimal.Decimal) error {
	u := map[string]interface{}{}
	u["balance"] = balance
	if util.IsDarwinia {
		u["kton_balance"] = additional
	}
	query := d.db.Model(account).Update(u)
	if query == nil || query.Error != nil || query.RowsAffected == 0 {
		return errors.New("update balance account error")
	}
	return nil
}

func (d *Dao) GetBalanceFromNetwork(c context.Context, address, module string) (decimal.Decimal, decimal.Decimal, error) {
	balanceModule := enableBalancesModule[module]
	if balanceModule == "" {
		return decimal.Zero, decimal.Zero, errors.New("disable module")
	}
	balance, additional, err := rpc.GetFreeBalance(nil, balanceModule, address, "")
	if err != nil {
		log.Error("GetBalanceFromNetwork error %v", err)
		return decimal.Zero, decimal.Zero, err
	}
	return balance.Div(decimal.New(1, int32(substrate.BalanceAccuracy))),
		additional.Div(decimal.New(1, int32(substrate.BalanceAccuracy))), nil
}

func (d *Dao) incrCountExtrinsic(c context.Context, account *model.ChainAccount, nonce int) bool {
	set := map[string]interface{}{
		"count_extrinsic": gorm.Expr("count_extrinsic + ?", 1),
		"nonce":           nonce,
	}
	if balance, _, err := d.GetBalanceFromNetwork(c, account.Address, BalanceModule); err == nil {
		set["balance"] = balance // cost gas update balance
	} else {
		log.Error("GetBalanceFromNetwork error %v", err)
	}
	query := d.db.Model(account).Update(set)
	if query != nil && query.RowsAffected > 0 {
		return true
	}
	return false
}

func (d *Dao) UpdateAccountCountExtrinsic(c context.Context, address string, nonce int) {
	account, err := d.TouchAccount(c, address)
	if err != nil {
		return
	}
	d.incrCountExtrinsic(c, account, nonce)
}

func (d *Dao) UpdateAccountBalance(c context.Context, account *model.ChainAccount, module string) (decimal.Decimal, decimal.Decimal, error) {
	balance, additional, err := d.GetBalanceFromNetwork(c, account.Address, module)
	if err == nil {
		_ = d.updateBalance(c, account, balance, additional)
	}
	return balance, additional, err
}

func (d *Dao) AccountAsJson(c context.Context, account *model.ChainAccount) *model.AccountJson {
	accountJson := model.AccountJson{
		AccountSampleJson: d.AccountSampleJson(c, account),
		AccountIndex:      ss58.EncodeAccountIndex(int64(account.AccountIndex), substrate.AddressType),
		Nonce:             account.Nonce,
	}
	if balance, additional, err := d.UpdateAccountBalance(c, account, BalanceModule); err == nil {
		accountJson.Balance = balance
		accountJson.KtonBalance = additional
	}
	return &accountJson
}

func (d *Dao) AccountSampleJson(c context.Context, account *model.ChainAccount) *model.AccountSampleJson {
	accuracy := int32(substrate.BalanceAccuracy)
	j := model.AccountSampleJson{
		Address:     ss58.Encode(account.Address, substrate.AddressType),
		Balance:     account.Balance,     // from chain
		KtonBalance: account.KtonBalance, // from chain
		RingLock:    account.RingLock.Div(decimal.New(1, accuracy)),
		KtonLock:    account.KtonLock.Div(decimal.New(1, accuracy)),
		Nickname:    account.Nickname,
	}
	return &j
}

func (d *Dao) ResetAccountNonce(c context.Context, address string) {
	account, err := d.TouchAccount(c, address)
	if err != nil {
		return
	}
	_ = d.db.Model(account).Update(model.ChainAccount{Nonce: 0})
}

func (d *Dao) UpdateAccountLock(c context.Context, address, currency string) error {
	balance, err := rpc.GetAccountLock(nil, address, currency)
	if err != nil {
		log.Error("UpdateAccountLock err %v", err)
		return err
	}
	u := map[string]interface{}{}
	if currency == Ring {
		u["ring_lock"] = balance
	} else if currency == Kton {
		u["kton_lock"] = balance
	}
	query := d.db.Model(model.ChainAccount{}).Where("address = ?", address).Update(u)
	if query == nil || query.Error != nil || query.RowsAffected == 0 {
		return errors.New("update balance lock error")
	}
	return nil
}

func (d *Dao) GetAccountList(c context.Context, page, row int, order, field string, queryWhere ...string) ([]*model.ChainAccount, int) {
	var accounts []*model.ChainAccount
	queryOrigin := d.db.Model(&model.ChainAccount{})
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

func (d *Dao) FindByIndex(c context.Context, index int) (*model.ChainAccount, error) {
	var account model.ChainAccount
	query := d.db.First(&account, &model.ChainAccount{AccountIndex: index})
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return nil, errors.New("touch account error")
	}
	return &account, nil
}
