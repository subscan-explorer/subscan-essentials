package service

import (
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockPluginStorageDao struct {
	mock.Mock
}

func (MockPluginStorageDao) FindBy(record interface{}, query interface{}, option *storage.Option) bool {
	return true
}

func (MockPluginStorageDao) AutoMigration(model interface{}) error {
	return nil
}

func (MockPluginStorageDao) AddIndex(model interface{}, indexName string, columns ...string) error {
	return nil
}

func (MockPluginStorageDao) AddUniqueIndex(model interface{}, indexName string, columns ...string) error {
	return nil
}

func (MockPluginStorageDao) Create(record interface{}) error {
	return nil
}

func (MockPluginStorageDao) Update(model interface{}, query interface{}, attr map[string]interface{}) error {
	return nil
}

func (MockPluginStorageDao) Delete(model interface{}, query interface{}) error {
	return nil
}

func (MockPluginStorageDao) SpecialMetadata(int) string {
	return ""
}

func (MockPluginStorageDao) RPCPool() *websocket.PoolConn {
	return nil
}

func (MockPluginStorageDao) SetPrefix(string) {
}

func TestService_pluginRegister(t *testing.T) {
	c := MockPluginStorageDao{}
	pluginRegister(c)
}

func Test_emitEvent(t *testing.T) {
	testSrv.emitEvent(&testBlock, &testEvent, decimal.Zero)
}

func Test_emitExtrinsic(t *testing.T) {
	testSrv.emitExtrinsic(&testBlock, &testSignedExtrinsic, []model.ChainEvent{testEvent})
}
