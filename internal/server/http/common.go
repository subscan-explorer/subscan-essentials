package http

import (
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"net/http"
	"time"
)

func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func now(c *bm.Context) {
	c.JSON(time.Now().Unix(), nil)
}

func systemStatus(c *bm.Context) {
	status := svc.DaemonHealth(c)
	c.JSON(status, nil)
}
