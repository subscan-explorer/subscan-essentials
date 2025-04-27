package dao

import (
	"context"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
)

type Account struct {
	Address    string `json:"address" gorm:"primaryKey;size:70"`
	EvmAccount string `json:"evm_account" gorm:"size:70;index:evm_account"`
}

func (a *Account) TableName() string {
	return "evm_accounts"
}

func TouchAccount(ctx context.Context, h160 string) error {
	if !address.VerifyEthereumAddress(h160) {
		return nil
	}
	addr := h160ToAccountIdByNetwork(ctx, h160, util.NetworkNode)
	if !sg.redis.SAdd(ctx, EvmAddressMemberKey, 86400, h160) {
		return nil
	}
	account := &Account{EvmAccount: h160, Address: addr}
	if err := sg.db.WithContext(ctx).Debug().Scopes(model.IgnoreDuplicate).Create(account).Error; err != nil {
		return err
	}
	return nil
}
