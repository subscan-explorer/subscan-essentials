package observer

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/itering/subscan/configs"
	"golang.org/x/exp/slog"
)

func newHealthCheckServer(c *configs.HealthCheck) *http.Server {
	if c == nil {
		panic("health check config is nil")
	}

	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.Default()
	e.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET"},
		MaxAge:          12 * time.Hour,
		ExposeHeaders:   []string{"Content-Length"},
		AllowHeaders:    []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
	}))

	health := func(c *gin.Context) {
		c.String(200, "ok")
	}
	e.GET("/", health)
	e.GET("/health", health)
	e.GET("/ping", health)

	srv := &http.Server{
		Addr:    c.Addr,
		Handler: e,
	}

	return srv
}

func startHealthCheckServer(srv *http.Server, stop chan struct{}) {
	go func() {
		slog.Info("Healthcheck server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Healthcheck server error", "error", err)
		}
	}()

	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server Shutdown failed", "error", err)
		}
		<-ctx.Done()
		slog.Info("Healthcheck server exiting")
	}()
}
