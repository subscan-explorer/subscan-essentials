package service

import (
	"context"
)

func (s *Service) UpdateChainMetadata(metadata map[string]interface{}) (err error) {
	c := context.TODO()
	err = s.Dao.SetMetadata(c, metadata)
	return
}

func (s *Service) GetCurrentBlockNum(c context.Context) (uint64, error) {
	return s.Dao.GetCurrentBlockNum(c)
}

func (s *Service) GetFinalizedBlockNum(c context.Context) (uint64, error) {
	return s.Dao.GetFinalizedBlockNum(c)
}
