package service

import (
	"github.com/itering/subscan-plugin/storage"
	internalDao "github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/model"
	"golang.org/x/exp/slog"
)

type Service struct {
	d  storage.Dao
	dd *internalDao.Dao
}

func New(d storage.Dao, dd *internalDao.Dao) *Service {
	return &Service{
		d:  d,
		dd: dd,
	}
}

func (s *Service) GetPayoutListJson(page, row int, address string) ([]model.Payout, int) {
	slog.Debug("GetPayoutListJson: ", page, row, address)
	return dao.GetPayoutList(s.d, page, row, address)
}
