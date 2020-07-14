package balance

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/http"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/shopspring/decimal"
)

var srv *service.Service

type Account struct {
	d storage.Dao
	e *bm.Engine
}

func New() *Account {
	return &Account{}
}

func (a *Account) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *Account) InitHttp(e *bm.Engine) {
	a.e = e
	http.Router(srv, a.e)
}

func (a *Account) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, events []storage.Event) error {
	return nil
}

func (a *Account) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	return nil
}

func (a *Account) Migrate() {
	db := a.d.DB()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.Account{},
	)
	db.Model(model.Account{}).AddUniqueIndex("address", "address")
}
