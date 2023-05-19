package service

import (
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/storage"
)

type Service struct {
	d storage.Dao
}

func (s *Service) Storage() storage.Dao {
	return s.d
}

func New(d storage.Dao) *Service {
	return &Service{
		d: d,
	}
}

func (s *Service) GetPayoutListJson(page, row int, address string, minEra uint32) ([]model.Payout, int) {
	return dao.GetPayoutList(s.d, page, row, address, minEra)
}

func (s *Service) GetRuntimeConstant(module, name string) *storage.RuntimeConstant {
	return s.d.GetRuntimeConstant(module, name)
}
