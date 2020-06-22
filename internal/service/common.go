package service

import (
	"context"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
)

func (s *Service) SetHeartBeat(action string) {
	ctx := context.TODO()
	s.dao.SetHeartBeatNow(ctx, action)
}

func (s *Service) GetSystemHeartBeat(ctx context.Context) map[string]bool {
	return s.dao.GetHeartBeatNow(ctx)
}

func (s *Service) Metadata() (map[string]string, error) {
	c := context.TODO()
	m, err := s.dao.GetMetadata(c)
	m["networkNode"] = util.NetworkNode
	m["blockTime"] = util.IntToString(substrate.BlockTime)
	m["commissionAccuracy"] = util.IntToString(substrate.CommissionAccuracy)
	m["balanceAccuracy"] = util.IntToString(substrate.BalanceAccuracy)
	m["addressType"] = util.IntToString(substrate.AddressType)
	return m, err
}
