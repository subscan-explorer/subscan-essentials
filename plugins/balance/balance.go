package balance

import (
	"context"
	"fmt"
	"github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/dao"
	"github.com/itering/subscan/plugins/balance/http"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
	"strings"
)

var srv *service.Service

type Balance struct {
	d storage.Dao
}

func (a *Balance) ConsumptionQueue() []string {
	return nil
}

func (a *Balance) Enable() bool {
	return true
}

func (a *Balance) ProcessBlock(context.Context, *storage.Block) error { return nil }

func (a *Balance) SetRedisPool(_ subscan_plugin.RedisPool) {}

func New() *Balance {
	return &Balance{}
}

func (a *Balance) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *Balance) InitHttp() []router.Http {
	return http.Router(srv)
}

func (a *Balance) ProcessExtrinsic(*storage.Block, *storage.Extrinsic, []storage.Event) error {
	return nil
}

func (a *Balance) ProcessEvent(_ *storage.Block, event *storage.Event, _ decimal.Decimal) error {
	if event == nil {
		return nil
	}
	var paramEvent []storage.EventParam
	util.UnmarshalAny(&paramEvent, event.Params)

	switch fmt.Sprintf("%s-%s", strings.ToLower(event.ModuleId), strings.ToLower(event.EventId)) {
	case strings.ToLower("System-NewAccount"):
		return dao.NewAccount(a.d, util.ToString(paramEvent[0].Value))
	}

	return nil
}

func (a *Balance) SubscribeExtrinsic() []string {
	return nil
}

func (a *Balance) SubscribeEvent() []string {
	return []string{"system"}
}

func (a *Balance) Version() string {
	return "0.1"
}

func (a *Balance) Migrate() {
	_ = a.d.AutoMigration(&model.Account{})
}

func (a *Balance) ExecWorker(context.Context, string, string, interface{}) error { return nil }
