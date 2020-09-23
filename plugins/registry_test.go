package plugins_test

import (
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TPlugin struct{}

func (a *TPlugin) InitDao(d storage.Dao) {}

func (a *TPlugin) InitHttp() []router.Http { return nil }

func (a *TPlugin) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, events []storage.Event) error {
	return nil
}

func (a *TPlugin) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	return nil
}

func (a *TPlugin) Migrate() {}

func (a *TPlugin) Version() string { return "0.1" }

func (a *TPlugin) SubscribeExtrinsic() []string { return nil }

func (a *TPlugin) SubscribeEvent() []string { return nil }

func TestRegister(t *testing.T) {
	plugins.Register("test", &TPlugin{})
	plugins.Register("test2", nil)
	plugins.Register("test", &TPlugin{})
	assert.NotNil(t, plugins.RegisteredPlugins["test"])
	assert.Nil(t, plugins.RegisteredPlugins["test2"])
}

func TestList(t *testing.T) {
	assert.Equal(t, len(plugins.List()), len(plugins.RegisteredPlugins))
}
