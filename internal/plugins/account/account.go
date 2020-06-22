package account

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/plugins/account/http"
	"github.com/itering/subscan/internal/plugins/account/service"
)

type Account struct {
	d *dao.Dao
	e *bm.Engine
}

func New() *Account {
	return &Account{}
}

func (a *Account) Init(d *dao.Dao, e *bm.Engine) error {
	a.d = d
	a.e = e
	return nil
}

func (a *Account) Http() error {
	srv := service.New(a.d)
	http.Router(srv, a.e)
	return nil
}

func (a *Account) ListModel() ([]interface{}, error) {
	return nil, nil
}

func (a *Account) ProcessExtrinsic() error {
	return nil
}

func (a *Account) ProcessEvent() error {
	return nil
}
