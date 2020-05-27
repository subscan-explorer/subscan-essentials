package http

import (
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/util/ss58"
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

func codecAddress(c *bm.Context) {
	address, _ := c.Params.Get("p")
	codec, _ := c.Params.Get("t")
	if codec == "encode" {
		address = ss58.Encode(address, substrate.AddressType)
	} else {
		address = ss58.Decode(address, substrate.AddressType)
	}
	c.JSON(address, nil)
}
