package service

import (
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/dao"
	"github.com/itering/subscan/plugins/balance/model"
)

type Service struct {
	d storage.Dao
}

func (s *Service) GetAccountListJson(page, row int, order string) ([]model.Account, int) {
	return dao.GetAccountList(s.d, page, row, order)
}

func New(d storage.Dao) *Service {
	return &Service{
		d: d,
	}
}
