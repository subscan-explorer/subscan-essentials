package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/plugins"
)

var (
	svc *service.Service
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *configs.Server, s *service.Service) *http.Server {
	var opts []http.ServerOption
	svc = s
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != "" {
		timeout, _ := time.ParseDuration(c.Http.Timeout)
		opts = append(opts, http.Timeout(timeout))
	}
	engine := http.NewServer(opts...)

	e := gin.New()
	e.Use(gin.Recovery())
	defer engine.HandlePrefix("/", e)
	initRouter(e)

	return engine
}

func initRouter(e *gin.Engine) {
	e.GET("ping", ping)
	// internal
	g := e.Group("/api")
	{
		g.POST("/now", now)
		s := g.Group("/scan")
		{
			s.POST("metadata", metadataHandle)

			// Block
			s.POST("blocks", blocksHandle)
			s.POST("block", blockHandle)

			// Extrinsic
			s.POST("extrinsics", extrinsicsHandle)
			s.POST("extrinsic", extrinsicHandle)
			// Event
			s.POST("events", eventsHandle)

			s.POST("check_hash", checkSearchHashHandle)

			// Runtime
			s.POST("runtime/metadata", runtimeMetadataHandle)
			s.POST("runtime/list", runtimeListHandler)

		}
		pluginRouter(g)
	}
}

func pluginRouter(g *gin.RouterGroup) {
	for name, plugin := range plugins.RegisteredPlugins {
		for _, r := range plugin.InitHttp() {
			g.Group("plugin").Group(name).POST(r.Router, func(context *gin.Context) {
				_ = r.Handle(context.Writer, context.Request)
			})
		}
	}
}
