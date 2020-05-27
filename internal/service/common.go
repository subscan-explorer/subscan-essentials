package service

import (
	"context"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/util"
)

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

func (s *Service) SetHeartBeat(action string) {
	ctx := context.TODO()
	s.dao.SetHeartBeatNow(ctx, action)
}

func (s *Service) GetSystemHeartBeat(ctx context.Context) map[string]bool {
	return s.dao.GetHeartBeatNow(ctx)
}

func (s *Service) GetTransactionCount(c context.Context) int {
	return s.dao.GetTransactionCount(c)
}

func (s *Service) Metadata() (map[string]string, error) {
	c := context.TODO()
	m, err := s.dao.GetMetadata(c)
	m["blockTime"] = util.IntToString(substrate.BlockTime)
	m["networkNode"] = util.NetworkNode
	m["commissionAccuracy"] = substrate.CommissionAccuracy
	m["addressType"] = util.IntToString(substrate.AddressType)
	return m, err
}

func (s *Service) UpdateAccountAllBalance(address string) {
	c := context.TODO()
	if account, err := s.dao.TouchAccount(c, address); err == nil {
		_, _, _ = s.dao.UpdateAccountBalance(c, account, "balances")
		_ = s.dao.UpdateAccountLock(c, account.Address, "ring")
	}
}
