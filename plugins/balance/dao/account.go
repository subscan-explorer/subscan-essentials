package dao

import (
	"context"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/model"
	bModel "github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAccountListCursor(db storage.DB, limit int, before, after *uint) ([]bModel.Account, bool, bool) {
	var accounts []bModel.Account
	d := db.GetDbInstance().(*gorm.DB)
	fetch := limit + 1
	var hasPrev, hasNext bool
	q := d.Model(bModel.Account{})
	if after != nil && *after > 0 {
		q = q.Where("id < ?", *after).Order("id desc")
	} else if before != nil && *before > 0 {
		q = q.Where("id > ?", *before).Order("id asc")
	} else {
		q = q.Order("id desc")
	}
	q = q.Limit(fetch).Find(&accounts)
	if q.Error != nil {
		return nil, false, false
	}
	if before != nil && *before > 0 {
		hasPrev = len(accounts) > limit
		if hasPrev {
			accounts = accounts[:limit]
		}
		for i, j := 0, len(accounts)-1; i < j; i, j = i+1, j-1 {
			accounts[i], accounts[j] = accounts[j], accounts[i]
		}
		hasNext = true
	} else {
		hasNext = len(accounts) > limit
		if hasNext {
			accounts = accounts[:limit]
		}
		hasPrev = after != nil && *after > 0
	}
	return accounts, hasPrev, hasNext
}

func GetAccountByAddress(ctx context.Context, db storage.DB, address string) *bModel.Account {
	var account bModel.Account
	d := db.GetDbInstance().(*gorm.DB)
	q := d.WithContext(ctx).Where("address = ?", address).First(&account)
	if q.Error != nil {
		return nil
	}
	return &account

}

func RefreshAccount(ctx context.Context, s *Storage, accountId string) error {
	accountId = address.Format(accountId)
	if accountId == "" {
		return nil
	}
	db := s.Dao.GetDbInstance().(*gorm.DB)
	var account = bModel.Account{Address: accountId}
	q := db.WithContext(ctx).Scopes(model.IgnoreDuplicate).Where("address = ?", accountId).FirstOrCreate(&account)
	if q.RowsAffected == 1 {
		_, _ = s.Pool.HINCRBY(ctx, model.MetadataCacheKey(), "total_account", 1)
	}
	return AfterAccountCreate(ctx, db, &account)
}

func AfterAccountCreate(ctx context.Context, db *gorm.DB, account *bModel.Account) error {
	accountDataRaw, err := rpc.ReadStorage(nil, "system", "account", "", account.Address)
	if err != nil {
		return err
	}
	accountData := new(bModel.AccountData)
	accountDataRaw.ToAny(accountData)
	return db.WithContext(ctx).Model(account).Where("address = ?", account.Address).UpdateColumns(map[string]interface{}{
		"nonce":    accountData.Nonce,
		"balance":  accountData.Data.Free.Add(accountData.Data.Reserved),
		"locked":   decimal.Max(accountData.Data.MiscFrozen, accountData.Data.FeeFrozen),
		"reserved": accountData.Data.Reserved,
	}).Error
}

func (s *Storage) AddOrUpdateItem(c context.Context, item interface{}, keys []string, updates ...string) *gorm.DB {
	var keyFields []clause.Column
	for _, key := range keys {
		keyFields = append(keyFields, clause.Column{Name: key})
	}
	db := s.Dao.GetDbInstance().(*gorm.DB)
	if len(updates) > 0 {
		return db.WithContext(c).Clauses(clause.OnConflict{
			Columns:   keyFields,
			DoUpdates: clause.AssignmentColumns(updates),
		}).Create(item)
	}
	return db.WithContext(c).Clauses(clause.OnConflict{
		Columns:   keyFields,
		UpdateAll: true,
	}).Create(item)
}
