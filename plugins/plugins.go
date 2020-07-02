package plugins

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/storage"
	"github.com/shopspring/decimal"
)

type Plugin interface {
	InitDao(d storage.Dao)

	InitHttp(e *bm.Engine)

	Http() error

	ProcessExtrinsic(int, *model.ChainExtrinsic, []model.ChainEvent) error

	ProcessEvent(spec, blockTimestamp int, blockHash string, event *model.ChainEvent, fee decimal.Decimal) error

	Migrate()
}
