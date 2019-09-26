package http

import (
	middlewares "subscan-end/internal/middleware"
	"subscan-end/internal/server/websocket"
	"subscan-end/internal/service"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var (
	svc   *service.Service
	WsHub *websocket.Hub
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
	initWs()
	initRouter(engine)
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine) {
	limiter := bm.NewRateLimiter(nil)
	e.Use(limiter.Limit())
	e.GET("socket", wsPullHandle) //Websocket
	g := e.Group("/api")
	{
		e.Ping(ping)
		g.GET("system/status", systemStatus)
		g.Use(middlewares.CORS())
		{
			g.POST("/now", now)
			s := g.Group("/scan")
			{
				s.POST("metadata", metadata)
				// Block
				s.POST("blocks", blocks)
				s.POST("block", block)
				// Extrinsics
				s.POST("extrinsics", extrinsics)
				s.POST("extrinsic", extrinsic)
				// Event
				s.POST("events", events)
				s.POST("event", event)
				//Search
				s.POST("search", search)
				s.POST("check_hash", checkSearchHash)
				// Log
				s.POST("logs", logs)
				s.POST("log", logInfo)
				// Transfer
				s.POST("transfers", transfers)
				// daily stat
				s.POST("daily", dailyStat)
			}
		}
	}
}

func initWs() {
	WsHub = websocket.NewHub()
	go WsHub.Run()
	websocket.NewMessageRouter(WsHub.Broadcast, svc)
}
