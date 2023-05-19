package http

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/itering/subscan/configs"
	middlewares "github.com/itering/subscan/internal/middleware"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/plugins"
)

var svc *service.ReadOnlyService

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *configs.Server, s *service.ReadOnlyService) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			tracing.Server(),
			metrics.Server(),
			validate.Validator(),
		),
	}

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
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(middlewares.CORS())
	defer engine.HandlePrefix("/", e)
	initRouter(e)
	return engine
}

func initRouter(e *gin.Engine) {
	e.GET("ping", ping)
	// internal
	g := e.Group("/api")
	{
		g.GET("system/status", systemStatus)
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

			s.POST("check_hash", checkSearchHash)

			// Runtime
			s.POST("runtime/metadata", runtimeMetadata)
			s.POST("runtime/list", runtimeList)

			// Plugin
			s.POST("plugins", pluginList)
			s.POST("plugins/ui", pluginUIConfig)
		}
		pluginRouter(g)
	}
}

func pluginRouter(g *gin.RouterGroup) {
	plug := g.Group("plugin")
	for name, plugin := range plugins.RegisteredPlugins {
		group := plug.Group(name)
		routers := plugin.InitHttp()
		for _, r := range routers {
			group.POST(r.Router, r.Handle)
		}
	}
}
