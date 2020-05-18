package service

import (
	"context"
	"subscan-end/libs/substrate"
	"subscan-end/utiles"
)

func (s *Service) UpdateChainMetadata(metadata map[string]interface{}) error {
	c := context.TODO()
	return s.dao.SetMetadata(c, metadata)
}

func (s *Service) GetChainMetadata() (map[string]string, error) {
	c := context.TODO()
	metadata, err := s.dao.GetMetadata(c)
	metadata["blockTime"] = utiles.IntToString(substrate.BlockTime)
	metadata["networkNode"] = utiles.NetworkNode
	return metadata, err
}
