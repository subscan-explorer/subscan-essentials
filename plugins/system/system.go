package system

import (
	"context"
	subscan_plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/system/model"
	"github.com/itering/subscan/plugins/system/service"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
)

var srv *service.Service

type System struct {
	d storage.Dao
}

func (a *System) ConsumptionQueue() []string {
	return nil
}

func (a *System) Enable() bool {
	return true
}

func (a *System) ProcessBlock(context.Context, *storage.Block) error {
	return nil
}

func New() *System {
	return &System{}
}

func (a *System) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *System) InitHttp() (routers []router.Http) {
	return nil
}

func (a *System) ProcessExtrinsic(*storage.Block, *storage.Extrinsic, []storage.Event) error {
	return nil
}

func (a *System) ProcessEvent(block *storage.Block, event *storage.Event, _ decimal.Decimal) error {
	var paramEvent []storage.EventParam
	_ = util.UnmarshalAny(&paramEvent, event.Params)
	switch event.EventId {
	case "ExtrinsicFailed":
		srv.ExtrinsicFailed(block.SpecVersion, event, paramEvent)
	}
	return nil
}

func (a *System) Migrate() {
	db := a.d
	_ = db.AutoMigration(&model.ExtrinsicError{})
	_ = db.AddUniqueIndex(&model.ExtrinsicError{}, "extrinsic_hash", "extrinsic_hash")
}

func (a *System) SetRedisPool(subscan_plugin.RedisPool) {

}

func (a *System) Version() string {
	return "0.1"
}

func (a *System) SubscribeExtrinsic() []string {
	return nil
}

func (a *System) SubscribeEvent() []string {
	return []string{"system"}
}

func (a *System) ExecWorker(context.Context, string, string, interface{}) error { return nil }
