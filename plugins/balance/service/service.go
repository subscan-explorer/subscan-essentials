package service

import (
	subscan_plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/dao"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/util/address"
)

type Service struct {
	d    storage.Dao
	pool subscan_plugin.RedisPool
}

func (s *Service) GetAccountListJson(page, row int) ([]model.Account, int) {
	list, count := dao.GetAccountList(s.d, page, row)
	for i := range list {
		list[i].Address = address.Format(list[i].Address)
	}
	return list, count
}

func New(d storage.Dao, pool subscan_plugin.RedisPool) *Service {
	return &Service{
		d:    d,
		pool: pool,
	}
}
