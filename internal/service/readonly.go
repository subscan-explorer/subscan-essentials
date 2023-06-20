package service

import (
	"sync"

	"github.com/itering/subscan/internal/dao"
)

type ReadOnlyService struct {
	dao dao.IReadOnlyDao

	metadataLock sync.RWMutex
}

func NewReadOnly() *ReadOnlyService {
	dao, dbStorage := dao.NewReadOnly()
	pluginRegister(dbStorage)
	return &ReadOnlyService{dao: dao}
}

func readOnlyWithDao(dao dao.IReadOnlyDao) *ReadOnlyService {
	return &ReadOnlyService{dao: dao}
}
