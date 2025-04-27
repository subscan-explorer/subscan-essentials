package evm

import (
	"context"
	"github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/evm/dao"
	"github.com/itering/subscan/plugins/evm/http"
	"github.com/itering/subscan/plugins/evm/workers"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type EVM struct {
	d      storage.Dao
	s      *dao.Storage
	enable bool
}

func (a *EVM) Enable() bool {
	return a.enable
}

func (a *EVM) ProcessBlock(ctx context.Context, block *storage.Block) error {
	return a.s.AddEvmBlock(ctx, uint(block.BlockNum), false)
}

func (a *EVM) SetRedisPool(pool subscan_plugin.RedisPool) {
	if a.Enable() {
		a.s = dao.Init(a.d.GetDbInstance().(*gorm.DB), pool)
	}
}

func New() *EVM {
	return &EVM{}
}

func (a *EVM) InitDao(d storage.Dao) {
	support := metadata.SupportModule()
	// check runtime module has  EVM or Revive module
	if !util.StringInSlice("EVM", support) && !util.StringInSlice("Revive", support) {
		util.Logger().Warning("EVM plugin is disabled because the runtime does not support EVM or Revive module")
		return
	}
	a.enable = true
	a.d = d
	a.Migrate()
}

func (a *EVM) InitHttp() []router.Http {
	return http.Router()
}

func (a *EVM) ProcessExtrinsic(*storage.Block, *storage.Extrinsic, []storage.Event) error {
	return nil
}

func (a *EVM) ProcessEvent(_ *storage.Block, _ *storage.Event, _ decimal.Decimal) error { return nil }

func (a *EVM) SubscribeExtrinsic() []string {
	return nil
}

func (a *EVM) SubscribeEvent() []string {
	return []string{"evm"}
}

func (a *EVM) Version() string {
	return "0.1"
}

func (a *EVM) ConsumptionQueue() []string {
	return []string{dao.Eip20Token, dao.Eip721Token, dao.Eip1155Token}
}

func (a *EVM) ExecWorker(ctx context.Context, queue, class string, raw interface{}) error {
	return workers.Emit(ctx, queue, class, raw)
}

func (a *EVM) Migrate() {
	for _, table := range a.s.Tables() {
		_ = a.d.AutoMigration(table)
	}
}
