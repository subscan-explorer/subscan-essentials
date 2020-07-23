package dao

import (
	"fmt"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/jinzhu/gorm"
)

func GetAccountList(db *gorm.DB, page, row int, order, field string, queryWhere ...string) ([]*model.Account, int) {
	var accounts []*model.Account
	queryOrigin := db.Model(&model.Account{})
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
