package balance

import (
	"context"
	"github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/dao"
	"github.com/itering/subscan/plugins/balance/http"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/shopspring/decimal"
	"strings"
)

var srv *service.Service

type Balance struct {
	d    storage.Dao
	pool subscan_plugin.RedisPool
}

func (a *Balance) ConsumptionQueue() []string {
	return nil
}

func (a *Balance) Enable() bool {
	return true
}

func (a *Balance) ProcessBlock(context.Context, *storage.Block) error { return nil }

func (a *Balance) SetRedisPool(pool subscan_plugin.RedisPool) {
	a.pool = pool
	srv = service.New(a.d, pool)
}

func New() *Balance {
	return &Balance{}
}

func (a *Balance) InitDao(d storage.Dao) {
	a.d = d
	a.Migrate()
}

func (a *Balance) InitHttp() []router.Http {
	return http.Router(srv)
}

func (a *Balance) ProcessExtrinsic(*storage.Block, *storage.Extrinsic, []storage.Event) error {
	return nil
}

func (a *Balance) ProcessEvent(block *storage.Block, event *storage.Event, _ decimal.Decimal) error {
	if event == nil {
		return nil
	}
	switch strings.ToLower(event.ModuleId) {
	case strings.ToLower("Balances"):
		return dao.EmitEvent(context.TODO(), &dao.Storage{Dao: a.d, Pool: a.pool}, event, block)
	}

	return nil
}

func (a *Balance) SubscribeExtrinsic() []string {
	return nil
}

func (a *Balance) SubscribeEvent() []string {
	return []string{"balances"}
}

func (a *Balance) Version() string {
	return "0.1"
}

func (a *Balance) Migrate() {
	_ = a.d.AutoMigration(&model.Account{})
	_ = a.d.AutoMigration(&model.Transfer{})
}

func (a *Balance) ExecWorker(context.Context, string, string, interface{}) error { return nil }

func (a *Balance) RefreshMetadata() {
	dao.RefreshMetadata(context.TODO(), &dao.Storage{Dao: a.d, Pool: a.pool})
}
