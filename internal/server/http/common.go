package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itering/subscan/util"
	"golang.org/x/exp/slog"
)

func ping(ctx *gin.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		slog.Warn("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func now(c *gin.Context) {
	util.ToJson(c, time.Now().Unix(), nil)
}

func systemStatus(c *gin.Context) {
	status := svc.DaemonHealth(c)
	util.ToJson(c, status, nil)
}
