package service

import (
	balance "github.com/itering/subscan/internal/plugins/balance/dao"
	"github.com/itering/subscan/internal/plugins/balance/model"
	"github.com/itering/subscan/internal/plugins/storage"
)

type Service struct {
	d storage.Dao
}

func (s *Service) GetAccountListJson(page, row int, order, field string, queryWhere ...string) ([]*model.Account, int) {
	return balance.GetAccountList(s.d.DB(), page, row, order, field, queryWhere...)
}

func New(d storage.Dao) *Service {
	return &Service{
		d: d,
	}
}
