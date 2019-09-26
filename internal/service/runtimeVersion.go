package service

import "context"

func (s *Service) CreateRuntimeVersion(name string, spec int) error {
	c := context.TODO()
	return s.dao.CreateRuntimeVersion(c, name, spec)
}
