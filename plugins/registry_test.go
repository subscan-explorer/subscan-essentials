package plugins

import (
	"testing"

	subscan_plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	internalDao "github.com/itering/subscan/internal/dao"
	scanModel "github.com/itering/subscan/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type TPlugin struct{}

func (a *TPlugin) InitDao(d storage.Dao, dd *internalDao.Dao) {}

func (a *TPlugin) InitHttp() []router.Http { return nil }

func (a *TPlugin) ProcessExtrinsic(block *scanModel.ChainBlock, extrinsic *scanModel.ChainExtrinsic, events []scanModel.ChainEvent) error {
	return nil
}

func (a *TPlugin) ProcessEvent(block *scanModel.ChainBlock, event *scanModel.ChainEvent, fee decimal.Decimal, extrinsic *scanModel.ChainExtrinsic) error {
	return nil
}

func (a *TPlugin) ProcessCall(block *scanModel.ChainBlock, call *scanModel.ChainCall, events []scanModel.ChainEvent, extrinsic *scanModel.ChainExtrinsic) error {
	return nil
}

func (a *TPlugin) Migrate() {}

func (a *TPlugin) Version() string { return "0.1" }

func (a *TPlugin) SubscribeExtrinsic() []string { return nil }

func (a *TPlugin) SubscribeEvent() []string { return nil }

func (a *TPlugin) SubscribeCall() []string { return nil }

func (a *TPlugin) UiConf() *subscan_plugin.UiConfig { return nil }

func TestRegister(t *testing.T) {
	register("test", &TPlugin{})
	register("test2", nil)
	register("test", &TPlugin{})
	assert.NotNil(t, RegisteredPlugins["test"])
	assert.Nil(t, RegisteredPlugins["test2"])
}

func TestList(t *testing.T) {
	assert.Equal(t, len(List()), len(RegisteredPlugins))
}
