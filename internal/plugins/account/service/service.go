package service

import (
	"github.com/itering/subscan/internal/dao"
	account "github.com/itering/subscan/internal/plugins/account/dao"
	"github.com/itering/subscan/internal/plugins/account/model"
)

type Service struct {
	d *dao.Dao
}

func (s *Service) GetAccountListJson(page, row int, order, field string, queryWhere ...string) ([]*model.ChainAccount, int) {
	return account.GetAccountList(s.d.Db, page, row, order, field, queryWhere...)
}

func New(d *dao.Dao) *Service {
	return &Service{
		d: d,
	}
}
