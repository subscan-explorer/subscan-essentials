package plugins_test

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TPlugin struct{}

func (a *TPlugin) InitDao(d storage.Dao) {}

func (a *TPlugin) InitHttp(e *bm.Engine) {}

func (a *TPlugin) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, events []storage.Event) error {
	return nil
}

func (a *TPlugin) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	return nil
}

func (a *TPlugin) Migrate() {}

func TestRegister(t *testing.T) {
	plugins.Register("test", &TPlugin{})
	assert.NotNil(t, plugins.RegisteredPlugins["test"])
}

func TestList(t *testing.T) {
	assert.Equal(t, len(plugins.List()), len(plugins.RegisteredPlugins))
}
