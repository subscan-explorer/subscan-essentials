package service

import (
	"context"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util"
)

func (s *Service) SetHeartBeat(action string) {
	ctx := context.TODO()
	s.Dao.SetHeartBeatNow(ctx, action)
}

func (s *Service) GetSystemHeartBeat(ctx context.Context) map[string]bool {
	return s.Dao.GetHeartBeatNow(ctx)
}

func (s *Service) GetTransactionCount(c context.Context) int {
	return s.Dao.GetTransactionCount(c)
}

func (s *Service) Metadata() (map[string]string, error) {
	c := context.TODO()
	m, err := s.Dao.GetMetadata(c)
	m["blockTime"] = util.IntToString(substrate.BlockTime)
	m["networkNode"] = util.NetworkNode
	m["commissionAccuracy"] = util.IntToString(substrate.CommissionAccuracy)
	m["addressType"] = util.IntToString(substrate.AddressType)
	return m, err
}
