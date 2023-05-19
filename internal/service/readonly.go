package service

import "github.com/itering/subscan/internal/dao"

type ReadOnlyService struct {
	dao dao.IReadOnlyDao
}

func NewReadOnly() *ReadOnlyService {
	dao, dbStorage := dao.NewReadOnly()
	pluginRegister(dbStorage)
	return &ReadOnlyService{dao: dao}
}

func readOnlyWithDao(dao dao.IReadOnlyDao) *ReadOnlyService {
	return &ReadOnlyService{dao: dao}
}
