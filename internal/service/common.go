package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/itering/subscan/util"
)

// Ping ping the resource.
func (s *ReadOnlyService) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

func (s *Service) SetHeartBeat(action string) {
	ctx := context.TODO()
	_ = s.dao.SetHeartBeatNow(ctx, action)
}

func (s *ReadOnlyService) DaemonHealth(ctx context.Context) map[string]bool {
	return s.dao.DaemonHealth(ctx)
}

func (s *ReadOnlyService) Metadata() (map[string]string, error) {
	c := context.TODO()
	m, err := s.dao.GetMetadata(c)
	m["networkNode"] = util.NetworkNode
	m["commissionAccuracy"] = util.CommissionAccuracy
	m["balanceAccuracy"] = util.BalanceAccuracy
	m["addressType"] = util.AddressType
	return m, err
}
