package service

import (
	"context"
	subscan_plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/storage"
	cmodel "github.com/itering/subscan/model"
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
		list[i].Address = address.Encode(list[i].Address)
	}
	return list, count
}

func (s *Service) GetAccountJson(ctx context.Context, addr string) *model.Account {
	account := dao.GetAccountByAddress(ctx, s.d, addr)
	if account == nil {
		return nil
	}
	account.Address = address.Encode(account.Address)
	return account
}

func (s *Service) GetTransferJson(ctx context.Context, addr string, blockNum uint, page, row int) ([]model.Transfer, int) {
	var opts []cmodel.Option
	if blockNum > 0 {
		opts = append(opts, cmodel.Where("block_num = ?", blockNum))
	}
	if addr != "" {
		opts = append(opts, cmodel.Where("sender = ? or receiver = ?", addr, addr))
	}
	opts = append(opts, cmodel.WithLimit(page*row, row))

	list, count := dao.Transfers(ctx, s.d, opts...)
	for index := range list {
		list[index].Sender = address.Encode(list[index].Sender)
		list[index].Receiver = address.Encode(list[index].Receiver)
	}
	return list, count
}

func New(d storage.Dao, pool subscan_plugin.RedisPool) *Service {
	return &Service{
		d:    d,
		pool: pool,
	}
}
