package system

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/model"
	system "github.com/itering/subscan/internal/plugins/system/model"
	"github.com/itering/subscan/internal/plugins/system/service"
	"github.com/itering/subscan/internal/util"
	"github.com/shopspring/decimal"
)

var srv *service.Service

type System struct {
	d *dao.Dao
	e *bm.Engine
}

func New() *System {
	return &System{}
}

func (a *System) InitDao(d *dao.Dao) {
	srv = service.New(a.d)
	a.d = d
	a.Migrate()
}
func (a *System) InitHttp(e *bm.Engine) {
	a.e = e
}

func (a *System) Http() error {
	return nil
}

func (a *System) ProcessExtrinsic(spec int, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) error {
	return nil
}

func (a *System) ProcessEvent(spec, blockTimestamp int, blockHash string, event *model.ChainEvent, fee decimal.Decimal) error {
	var paramEvent []model.EventParam
	util.UnmarshalToAnything(&paramEvent, event.Params)
	switch event.EventId {
	case "ExtrinsicFailed":
		srv.ExtrinsicFailed(spec, blockTimestamp, blockHash, event, paramEvent)
	}
	return nil
}

func (a *System) Migrate() {
	db := a.d.Db
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&system.ExtrinsicError{},
	)
	db.Model(system.ExtrinsicError{}).AddUniqueIndex("extrinsic_hash", "extrinsic_hash")
}
