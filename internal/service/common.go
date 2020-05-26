package service

import (
	"context"
	"github.com/itering/subscan/libs/substrate"
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
