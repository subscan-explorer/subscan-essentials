package di

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	"github.com/itering/subscan/internal/service"
)

//go:generate kratos tool wire
type App struct {
	svc  *service.Service
	http *bm.Engine
}

func NewApp(
	svc *service.Service, h *bm.Engine, g *warden.Server,
) (
	app *App, closeFunc func(), err error,
) {
	app = &App{
		svc:  svc,
		http: h,
	}
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err := h.Shutdown(ctx); err != nil {
			log.Error("httpSrv.Shutdown error(%v)", err)
		}
		cancel()
	}
	return
}
