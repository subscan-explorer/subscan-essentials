package http

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	"github.com/itering/subscan/plugins/balance/service"
)

var (
	svc *service.Service
)

func Router(s *service.Service, e *bm.Engine) {
	svc = s
	g := e.Group("/api")
	{
		s := g.Group("/scan")
		{
			s.POST("accounts", accounts)
		}
	}
}

func accounts(c *bm.Context) {
	p := new(struct {
		Row        int    `json:"row" validate:"min=1,max=100"`
		Page       int    `json:"page" validate:"min=0"`
		Order      string `json:"order" validate:"omitempty,oneof=desc asc"`
		OrderField string `json:"order_field" validate:"omitempty"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	var query []string
	list, count := svc.GetAccountListJson(p.Page, p.Row, p.Order, p.OrderField, query...)
	c.JSON(map[string]interface{}{
		"list": list, "count": count,
	}, nil)
}
