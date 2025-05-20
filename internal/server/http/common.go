package http

import (
	"net/http"
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Printf("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func now(c *gin.Context) {
	toJson(c, time.Now().Unix(), nil)
}
