package service

import (
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"subscan-end/internal/dao"
	"subscan-end/internal/model"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao *dao.Dao
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}
	model.InitEsClient()
	s.AfterInit()
	return s
}

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

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) AfterInit() {
	s.dao.Migration()
}
