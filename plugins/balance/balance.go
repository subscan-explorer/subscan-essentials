package balance

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/balance/http"
	model2 "github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/itering/subscan/plugins/storage"
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
	srv = service.New(a.d)
	a.d = d
	a.Migrate()
}
func (a *Account) InitHttp(e *bm.Engine) {
	a.e = e
}

func (a *Account) Http() error {
	http.Router(srv, a.e)
	return nil
}

func (a *Account) ProcessExtrinsic(spec int, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) error {
	return nil
}

func (a *Account) ProcessEvent(spec, blockTimestamp int, blockHash string, event *model.ChainEvent, fee decimal.Decimal) error {
	return nil
}

func (a *Account) Migrate() {
	db := a.d.DB()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model2.Account{},
	)
	db.Model(model2.Account{}).AddUniqueIndex("address", "address")
}
