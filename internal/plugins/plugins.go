package plugins

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/dao"
)

type Plugin interface {
	Init(d *dao.Dao, e *bm.Engine) error

	Http() error

	ListModel() ([]interface{}, error)

	ProcessExtrinsic() error

	ProcessEvent() error
}
