package service

import (
	"github.com/itering/subscan/internal/dao"
	balance "github.com/itering/subscan/internal/plugins/balance/dao"
	"github.com/itering/subscan/internal/plugins/balance/model"
)

type Service struct {
	d *dao.Dao
}

func (s *Service) GetAccountListJson(page, row int, order, field string, queryWhere ...string) ([]*model.ChainAccount, int) {
	return balance.GetAccountList(s.d.Db, page, row, order, field, queryWhere...)
}

func New(d *dao.Dao) *Service {
	return &Service{
		d: d,
	}
}
