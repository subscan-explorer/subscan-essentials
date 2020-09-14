package balance

import (
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/http"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/shopspring/decimal"
)

var srv *service.Service

type Balance struct {
	d storage.Dao
}

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

func (a *Balance) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, events []storage.Event) error {
	return nil
}

func (a *Balance) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	return nil
}

func (a *Balance) SubscribeExtrinsic() []string {
	return nil
}

func (a *Balance) SubscribeEvent() []string {
	return nil
}

func (a *Balance) Version() string {
	return "0.1"
}

func (a *Balance) Migrate() {
	_ = a.d.AutoMigration(&model.Account{})
	_ = a.d.AddUniqueIndex(&model.Account{}, "address", "address")
}
