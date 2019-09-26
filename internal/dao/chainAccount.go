package dao

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"subscan-end/internal/model"
	"subscan-end/utiles"
	"subscan-end/utiles/ss58"
)

func (d *Dao) TouchAccount(c context.Context, address string) (*model.ChainAccount, error) {
	var account model.ChainAccount
	query := d.db.FirstOrCreate(&account, &model.ChainAccount{Address: address})
	if query == nil || query.Error != nil {
		return nil, errors.New("touch account error")
	}
	if query.RowsAffected > 0 { //create
		_ = d.IncrMetadata(c, "count_account", 1)
	}
	return &account, nil
}

func (d *Dao) UpdateBalance(c context.Context, account *model.ChainAccount, balance decimal.Decimal, module string, countExtrinsic int) error {
	u := map[string]interface{}{"count_extrinsic": gorm.Expr("count_extrinsic + ?", countExtrinsic)}
	if module == "balances" {
		u["balance"] = balance
	} else if module == "kton" {
		u["kton_balance"] = balance
	}
	query := d.db.Model(account).Update(u)
	if query == nil || query.Error != nil || query.RowsAffected == 0 {
		return errors.New("update balance account error")
	}
	return nil
}

func (d *Dao) GetBalanceFromNetwork(c context.Context, address, module string) (decimal.Decimal, error) {
	balanceModule := enableBalancesModule[module]
	if balanceModule == "" {
		return decimal.Zero, errors.New("disable module")
	}
	if balance, err := d.substrateApi.GetFreeBalance(balanceModule, address); err != nil {
		return decimal.Zero, err
	} else {
		return utiles.BigToDecimal(balance), nil
	}
}

func (d *Dao) GetBalanceAndUpdate(c context.Context, account *model.ChainAccount, module string) error {
	if account == nil {
		return nil
	}
	if balance, err := d.GetBalanceFromNetwork(c, account.Address, module); err == nil {
		return d.UpdateBalance(c, account, balance, module, 1)
	}
	return nil
}

func (d *Dao) IncrCountExtrinsic(c context.Context, account *model.ChainAccount, countExtrinsic int) {
	d.db.Model(account).Update(map[string]interface{}{"count_extrinsic": gorm.Expr("count_extrinsic + ?", countExtrinsic)})
}

func (d *Dao) UpdateAccountCountExtrinsic(c context.Context, address string) {
	account, err := d.TouchAccount(c, address)
	if err != nil {
		return
	}
	d.IncrCountExtrinsic(c, account, 1)
}

func (d *Dao) UpdateAccountBalance(c context.Context, address, module string) {
	account, err := d.TouchAccount(c, address)
	if err != nil {
		return
	}
	if balance, err := d.GetBalanceFromNetwork(c, account.Address, module); err == nil {
		_ = d.UpdateBalance(c, account, balance, module, 0)
	}
}

func (d *Dao) UpdateAccountIndex(c context.Context, address string, accountIndex int) {
	account, err := d.TouchAccount(c, address)
	if err != nil {
		return
	}
	_ = d.db.Model(account).Update(model.ChainAccount{AccountIndex: accountIndex})
}

func (d *Dao) AccountAsJson(c context.Context, account *model.ChainAccount) *model.AccountJson {
	accountJson := model.AccountJson{
		Address:        ss58.Encode(account.Address),
		Balance:        account.Balance,
		KtonBalance:    account.KtonBalance,
		CountExtrinsic: account.CountExtrinsic,
		Nonce:          account.CountExtrinsic,
		AccountIndex:   account.AccountIndex,
	}
	return &accountJson
}
