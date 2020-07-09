package di

import (
	"context"
	"github.com/itering/subscan/internal/server/http"
	"github.com/itering/subscan/plugins"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/service"
)

type App struct {
	svc  *service.Service
	http *bm.Engine
}

func InitApp() (*App, func(), error) {
	serviceService := service.New()
	engine := http.New(serviceService)
	app, cleanup, err := newApp(serviceService, engine)

	// load plugins
	for _, plugin := range plugins.RegisteredPlugins {
		plugin.InitHttp(engine)
	}

	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		cleanup()
	}, nil
}

func newApp(
	svc *service.Service, h *bm.Engine,
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
		svc.Close()
	}
	return
}
