package service

import (
	"context"
)

func (s *Service) UpdateChainMetadata(metadata map[string]interface{}) (err error) {
	c := context.TODO()
	err = s.dao.SetMetadata(c, metadata)
	return
}

func (s *Service) GetCurrentBlockNum(c context.Context) (uint64, error) {
	return s.dao.GetCurrentBlockNum(c)
}

func (s *Service) GetFinalizedBlockNum(c context.Context) (uint64, error) {
	return s.dao.GetFinalizedBlockNum(c)
}
