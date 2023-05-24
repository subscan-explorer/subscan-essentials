package http

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/plugins"
)

var svc *service.ReadOnlyService

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *configs.Server, s *service.ReadOnlyService) *gin.Engine {
	svc = s

	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.Logger())
	e.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		MaxAge:          12 * time.Hour,
		ExposeHeaders:   []string{"Content-Length"},
		AllowHeaders:    []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
	}))
	initRouter(e)
	return e
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
