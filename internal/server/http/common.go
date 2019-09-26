package http

import (
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"net/http"
	"time"
)

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func now(c *bm.Context) {
	c.JSON(time.Now().Unix(), nil)
}

func systemStatus(c *bm.Context) {
	status := svc.GetSystemHeartBeat(c)
	c.JSON(status, nil)
}
