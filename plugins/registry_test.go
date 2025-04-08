package plugins

import (
	"context"
	"github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TPlugin struct{}

func (a *TPlugin) ProcessBlock(ctx context.Context, block *storage.Block) error {
	return nil
}

func (a *TPlugin) SetRedisPool(pool subscan_plugin.RedisPool) {
}

func (a *TPlugin) Enable() bool {
	return true
}

func (a *TPlugin) ConsumptionQueue() []string {
	// TODO implement me
	panic("implement me")
}

func (a *TPlugin) ExecWorker(ctx context.Context, queue, class string, raw interface{}) error {
	// TODO implement me
	panic("implement me")
}

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
	register("test", &TPlugin{})
	register("test2", nil)
	register("test", &TPlugin{})
	assert.NotNil(t, RegisteredPlugins["test"])
	assert.Nil(t, RegisteredPlugins["test2"])
}
