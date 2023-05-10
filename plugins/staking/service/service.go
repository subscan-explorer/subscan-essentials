package service

import (
	internalDao "github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/storage"
)

type Service struct {
	d  storage.Dao
	Dd *internalDao.Dao
}

func (s *Service) Storage() storage.Dao {
	return s.d
}

func (s *Service) Dao() *internalDao.Dao {
	return s.Dd
}

func New(d storage.Dao, dd *internalDao.Dao) *Service {
	return &Service{
		d:  d,
		Dd: dd,
	}
}

func (s *Service) GetPayoutListJson(page, row int, address string) ([]model.Payout, int) {
	res, count := dao.GetPayoutList(s.d, page, row, address)
	return res, count
}
