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

func (s *Service) GetAccountListCursor(_ context.Context, limit int, before, after *uint) ([]model.Account, map[string]interface{}) {
	list, hasPrev, hasNext := dao.GetAccountListCursor(s.d, limit, before, after)
	for i := range list {
		list[i].Address = address.Encode(list[i].Address)
	}
	var start, end *uint
	if len(list) > 0 {
		start = &list[0].ID
		end = &list[len(list)-1].ID
	}
	return list, map[string]interface{}{
		"start_cursor":      start,
		"end_cursor":        end,
		"has_previous_page": hasPrev,
		"has_next_page":     hasNext,
	}
}

func (s *Service) GetAccountJson(ctx context.Context, addr string) *model.Account {
	account := dao.GetAccountByAddress(ctx, s.d, addr)
	if account == nil {
		return nil
	}
	account.Address = address.Encode(account.Address)
	return account
}

func (s *Service) GetTransferCursor(ctx context.Context, addr string, blockNum uint, limit int, before, after *uint) ([]model.Transfer, map[string]interface{}) {
	var opts []cmodel.Option
	if blockNum > 0 {
		opts = append(opts, cmodel.Where("block_num = ?", blockNum))
	}
	if addr != "" {
		opts = append(opts, cmodel.Where("sender = ? or receiver = ?", addr, addr))
	}
	list, hasPrev, hasNext := dao.TransfersCursor(ctx, s.d, limit, before, after, opts...)
	for index := range list {
		list[index].Sender = address.Encode(list[index].Sender)
		list[index].Receiver = address.Encode(list[index].Receiver)
	}
	var start, end *uint
	if len(list) > 0 {
		start = &list[0].Id
		end = &list[len(list)-1].Id
	}
	return list, map[string]interface{}{
		"start_cursor":      start,
		"end_cursor":        end,
		"has_previous_page": hasPrev,
		"has_next_page":     hasNext,
	}
}

func New(d storage.Dao, pool subscan_plugin.RedisPool) *Service {
	return &Service{
		d:    d,
		pool: pool,
	}
}
