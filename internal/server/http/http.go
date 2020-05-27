package http

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/middleware"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/internal/service/scan"
)

var (
	svc *service.Service
	ss  *scan.Service
)

func New(s *service.Service) (engine *bm.Engine) {
	var (
		hc struct {
			Server *bm.ServerConfig
		}
	)
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	svc = s
	engine = bm.DefaultServer(hc.Server)
	engine.HandleMethodNotAllowed = false
	initRouter(engine)
	ss = svc.NewScan()
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine) {
	limiter := bm.NewRateLimiter(nil)
	e.Use(limiter.Limit(), middlewares.CORS())

	e.Ping(ping)
	// internal
	g := e.Group("/api")
	{
		g.GET("system/status", systemStatus)
		g.GET("tools/ss58/:p/:t", codecAddress)
		g.POST("/now", now)
		s := g.Group("/scan")
		{
			s.POST("metadata", metadata)
			// Block
			s.POST("blocks", blocks)
			s.POST("block", block)

			// Extrinsic
			s.POST("extrinsics", extrinsics)
			s.POST("extrinsic", extrinsic)
			// Event
			s.POST("events", events)
			s.POST("event", event)
			// Search
			s.POST("search", search)
			s.POST("check_hash", checkSearchHash)
			// Log
			s.POST("logs", logs)
			s.POST("log", logInfo)

			s.POST("accounts", accounts)

			s.POST("runtime/metadata", runtimeMetadata)
			s.POST("runtime/list", runtimeList)

		}
	}

}
